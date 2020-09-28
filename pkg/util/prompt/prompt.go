package prompt

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/p9c/pkg/app/slog"
	"os"
	"strings"

	"github.com/btcsuite/golangcrypto/ssh/terminal"

	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/util/legacy/keystore"
)

// ProvideSeed is used to prompt for the wallet seed which maybe required during
// upgrades.
func ProvideSeed() (seed []byte, err error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter existing wallet seed: ")
		var seedStr string
		if seedStr, err = reader.ReadString('\n'); slog.Check(err) {
			return
		}
		seedStr = strings.TrimSpace(strings.ToLower(seedStr))
		if seed, err = hex.DecodeString(seedStr); slog.Check(err) ||
			len(seed) < hdkeychain.MinSeedBytes ||
			len(seed) > hdkeychain.MaxSeedBytes {
			slog.Errorf("Invalid seed specified.  Must be a hexadecimal value that is at least %d bits and at"+
				" most %d bits", hdkeychain.MinSeedBytes*8, hdkeychain.MaxSeedBytes*8)
			continue
		}
		return
	}
}

// ProvidePrivPassphrase is used to prompt for the private passphrase which
// maybe required during upgrades.
func ProvidePrivPassphrase() (pass []byte, err error) {
	prompt := "enter the private passphrase for your wallet: "
	for {
		fmt.Print(prompt)
		if pass, err = terminal.ReadPassword(int(os.Stdin.Fd())); slog.Check(err) {
			return
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}
		return
	}
}

// promptList prompts the user with the given prefix, list of valid responses,
// and default list entry to use.  The function will repeat the prompt to the
// user until they enter a valid response.
func promptList(reader *bufio.Reader, prefix string, validResponses []string, defaultEntry string) (reply string, err error) {
	var prompt string
	// Setup the prompt according to the parameters.
	validStrings := strings.Join(validResponses, "/")
	if defaultEntry != "" {
		prompt = fmt.Sprintf("%s (%s) [%s]: ", prefix, validStrings, defaultEntry)
	} else {
		prompt = fmt.Sprintf("%s (%s): ", prefix, validStrings)
	}
	// Prompt the user until one of the valid responses is given.
	for {
		fmt.Print(prompt)
		if reply, err = reader.ReadString('\n'); slog.Check(err) {
			return
		}
		reply = strings.TrimSpace(strings.ToLower(reply))
		if reply == "" {
			reply = defaultEntry
		}
		for _, validResponse := range validResponses {
			if reply == validResponse {
				return
			}
		}
	}
}

// promptListBool prompts the user for a boolean (yes/no) with the given prefix. The function will repeat the prompt to
// the user until they enter a valid response.
func promptListBool(reader *bufio.Reader, prefix string, defaultEntry string) (resp bool, err error) {
	// Setup the valid responses.
	valid := []string{"n", "no", "y", "yes"}
	var response string
	if response, err = promptList(reader, prefix, valid, defaultEntry); slog.Check(err) {
		return
	}
	return response == "yes" || response == "y", nil
}

// promptPass prompts the user for a passphrase with the given prefix.  The
// function will ask the user to confirm the passphrase and will repeat the
// prompts until they enter a matching response.
func promptPass(reader *bufio.Reader, prefix string, confirm bool) (pass []byte, err error) {
	// Prompt the user until they enter a passphrase.
	prompt := fmt.Sprintf("%s: ", prefix)
	for {
		fmt.Print(prompt)
		if pass, err = terminal.ReadPassword(int(os.Stdin.Fd())); slog.Check(err) {
			return
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}
		if !confirm {
			return
		}
		fmt.Print("Confirm passphrase: ")
		var pass2 []byte
		if pass2, err = terminal.ReadPassword(int(os.Stdin.Fd())); slog.Check(err) {
			return
		}
		fmt.Print("\n")
		pass2 = bytes.TrimSpace(pass2)
		if !bytes.Equal(pass, pass2) {
			slog.Error("The entered passphrases do not match")
			continue
		}
		return
	}
}

// PrivatePass prompts the user for a private passphrase with varying behavior
// depending on whether the passed legacy keystore exists.  When it does, the
// user is prompted for the existing passphrase which is then used to unlock it.
// On the other hand, when the legacy keystore is nil, the user is prompted for
// a new private passphrase.  All prompts are repeated until the user enters a
// valid response.
func PrivatePass(reader *bufio.Reader, legacyKeyStore *keystore.Store) (privPass []byte, err error) {
	// When there is not an existing legacy wallet, simply prompt the user
	// for a new private passphase and return it.
	if legacyKeyStore == nil {
		return promptPass(reader,
			"Creating new wallet\n\nEnter the private passphrase for your new wallet", true)
	}
	// At this point, there is an existing legacy wallet, so prompt the user
	// for the existing private passphrase and ensure it properly unlocks
	// the legacy wallet so all of the addresses can later be imported.
	slog.Warn("You have an existing legacy wallet.  All addresses from your existing legacy wallet will be " +
		"imported into the new wallet format.")
	for {
		if privPass, err = promptPass(reader, ""+
			"Enter the private passphrase for your existing wallet", false); slog.Check(err) {
			return
		}
		// Keep prompting the user until the passphrase is correct.
		if err = legacyKeyStore.Unlock([]byte(privPass)); slog.Check(err) {
			if err == keystore.ErrWrongPassphrase {
				slog.Error(err)
				continue
			}
			return
		}
		return
	}
}

// PublicPass prompts the user whether they want to add an additional layer of
// encryption to the wallet.  When the user answers yes and there is already a
// public passphrase provided via the passed config, it prompts them whether or
// not to use that configured passphrase.  It will also detect when the same
// passphrase is used for the private and public passphrase and prompt the user
// if they are sure they want to use the same passphrase for both.  Finally, all
// prompts are repeated until the user enters a valid response.
func PublicPass(reader *bufio.Reader, privPass []byte, defaultPubPassphrase, configPubPassphrase []byte,
) (pubPass []byte, err error) {
	pubPass = defaultPubPassphrase
	var usePubPass bool
	if usePubPass, err = promptListBool(reader,
		"Do you want to add an additional layer of encryption for public data?", "no"); slog.Check(err) {
		return
	}
	if !usePubPass {
		return
	}
	if !bytes.Equal(configPubPassphrase, pubPass) {
		var useExisting bool
		if useExisting, err = promptListBool(reader, "Use the existing configured public passphrase for "+
			"encryption of public data?", "no"); slog.Check(err) {
			return
		}
		if useExisting {
			pubPass = configPubPassphrase
			return
		}
	}
	for {
		if pubPass, err = promptPass(reader,
			"Enter the public passphrase for your new wallet", true); err != nil {
			slog.Error(err)
			return
		}
		if bytes.Equal(pubPass, privPass) {
			var useSamePass bool
			if useSamePass, err = promptListBool(reader,
				"Are you sure want to use the same passphrase for public and private data?", "no"); slog.Check(err) {
				return
			}
			if useSamePass {
				break
			}
			continue
		}
		break
	}
	fmt.Println("NOTE: Use the --walletpass option to configure your public passphrase.")
	return
}

// Seed prompts the user whether they want to use an existing wallet generation
// seed.  When the user answers no, a seed will be generated and displayed to
// the user along with prompting them for confirmation.  When the user answers
// yes, a the user is prompted for it.  All prompts are repeated until the user
// enters a valid response.
func Seed(reader *bufio.Reader) (seed []byte, err error) {
	// Ascertain the wallet generation seed.
	var useUserSeed bool
	if useUserSeed, err = promptListBool(reader,
		"Do you have an existing wallet seed you want to use?", "no"); slog.Check(err) {
		return
	}
	if !useUserSeed {
		if seed, err = hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen); slog.Check(err) {
			return
		}
		fmt.Println("\nYour wallet generation seed is:")
		fmt.Printf("\n%x\n\n", seed)
		fmt.Print("IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet " +
			"without it.\n\n")
		fmt.Print("Please keep in mind that anyone who has access to the seed can also restore your wallet " +
			"thereby giving them access to all your funds, so it is imperative that you keep it in a " +
			"secure location.\n\n")
		for {
			fmt.Print(`Once you have stored the seed in a safe and secure location, enter "OK" to continue: `)
			var confirmSeed string
			if confirmSeed, err = reader.ReadString('\n'); slog.Check(err) {
				return
			}
			confirmSeed = strings.TrimSpace(confirmSeed)
			confirmSeed = strings.Trim(confirmSeed, `"`)
			if confirmSeed == "OK" {
				break
			}
		}
		return
	}
	for {
		fmt.Print("Enter existing wallet seed: ")
		var seedStr string
		if seedStr, err = reader.ReadString('\n'); slog.Check(err) {
			return
		}
		seedStr = strings.TrimSpace(strings.ToLower(seedStr))
		seed, err = hex.DecodeString(seedStr)
		if err != nil || len(seed) < hdkeychain.MinSeedBytes || len(seed) > hdkeychain.MaxSeedBytes {
			slog.Errorf("Invalid seed specified.  Must be a hexadecimal value that is at least %d bits and at "+
				"most %d bits\n", hdkeychain.MinSeedBytes*8, hdkeychain.MaxSeedBytes*8)
			continue
		}
		return
	}
}
