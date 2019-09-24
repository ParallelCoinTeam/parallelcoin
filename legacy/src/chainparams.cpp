// Copyright (c) 2010 Satoshi Nakamoto
// Copyright (c) 2009-2012 The Bitcoin developers
// Distributed under the MIT/X11 software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include "assert.h"

#include "chainparams.h"
#include "core.h"
#include "protocol.h"
#include "util.h"

//
// Main network
//

unsigned int pnSeed[] =
{
    
};

class CMainParams : public CChainParams {
public:
    CMainParams() {
        // The message start string is designed to be unlikely to occur in normal data.
        pchMessageStart[0] = 0xcd;
        pchMessageStart[1] = 0x08;
        pchMessageStart[2] = 0xac;
        pchMessageStart[3] = 0xff;
        vAlertPubKey = ParseHex("04fe8a17065a19b88f150305dda9ad209acb6bbe9b07cfaa36bcb383630b5b6f4c406a03acd0d77e2cbc4e9093c327c4f41df01c5791458066da7c1c337d0d15d3");
        nDefaultPort = 11047;
        nRPCPort = 11048;
        bnProofOfWorkLimit[ALGO_SHA256D] = CBigNum(~uint256(0) >> 20);
        bnProofOfWorkLimit[ALGO_SCRYPT]  = CBigNum(~uint256(0) >> 20);
        nSubsidyHalvingInterval = 250000;
  
        const char* pszTimestamp = "NYTimes 2014-07-19 - Dell Begins Accepting Bitcoin";
        CTransaction txNew;
        txNew.vin.resize(1);
        txNew.vout.resize(1);
        txNew.vin[0].scriptSig = CScript() << 486604799 << CBigNum(4) << vector<unsigned char>((const unsigned char*)pszTimestamp, (const unsigned char*)pszTimestamp + strlen(pszTimestamp));
        txNew.vout[0].nValue = 1000 * COIN;
        txNew.vout[0].scriptPubKey = CScript() << ParseHex("04e0d27172510c6806889740edafe6e63eb23fca32786fccfdb282bb2876a9f43b228245df057661ff943f6150716a20ea1851e8a7e9f54e620297664618438dae") << OP_CHECKSIG;
        genesis.vtx.push_back(txNew);
        genesis.hashPrevBlock = 0;
        genesis.hashMerkleRoot = genesis.BuildMerkleTree();
        genesis.nVersion = BLOCK_VERSION_DEFAULT;
        genesis.nTime    = 1405742300;
        genesis.nBits    = 0x1e0fffff;
        genesis.nNonce   = 66177;

		hashGenesisBlock = genesis.GetHash();
        
        assert(hashGenesisBlock == uint256("0x000009f0fcbad3aac904d3660cfdcf238bf298cfe73adf1d39d14fc5c740ccc7"));
        assert(genesis.hashMerkleRoot == uint256("0xc843eae4658e3a51d2f280c36376ce56dc71a6c70e4b1c5ad2d7a9316f9b9ab7"));

        vSeeds.push_back(CDNSSeedData("seed1.parallelcoin.info", "seed1.parallelcoin.info"));
        vSeeds.push_back(CDNSSeedData("seed3.parallelcoin.info", "seed3.parallelcoin.info"));
        vSeeds.push_back(CDNSSeedData("seed2.parallelcoin.info", "seed2.parallelcoin.info"));
        vSeeds.push_back(CDNSSeedData("seed4.parallelcoin.info", "seed4.parallelcoin.info"));
        vSeeds.push_back(CDNSSeedData("seed5.parallelcoin.info", "seed5.parallelcoin.info"));

        base58Prefixes[PUBKEY_ADDRESS] = 83;
        base58Prefixes[SCRIPT_ADDRESS] = 9;
        base58Prefixes[SECRET_KEY] = 178;

        // Convert the pnSeeds array into usable address objects.
        for (unsigned int i = 0; i < ARRAYLEN(pnSeed); i++)
        {
            // It'll only connect to one or two seed nodes because once it connects,
            // it'll get a pile of addresses with newer timestamps.
            // Seed nodes are given a random 'last seen time' 
            const int64 nTwoDays = 2 * 24 * 60 * 60;
            struct in_addr ip;
            memcpy(&ip, &pnSeed[i], sizeof(ip));
            CAddress addr(CService(ip, GetDefaultPort()));
            addr.nTime = GetTime() - GetRand(nTwoDays) - nTwoDays;
            vFixedSeeds.push_back(addr);
        }
    }

    virtual const CBlock& GenesisBlock() const { return genesis; }
    virtual Network NetworkID() const { return CChainParams::MAIN; }

    virtual const vector<CAddress>& FixedSeeds() const {
        return vFixedSeeds;
    }
protected:
    CBlock genesis;
    vector<CAddress> vFixedSeeds;
};
static CMainParams mainParams;


//
// Testnet (v3)
//
class CTestNetParams : public CMainParams {
public:
    CTestNetParams() {
        // The message start string is designed to be unlikely to occur in normal data.
        pchMessageStart[0] = 0x08;
        pchMessageStart[1] = 0xb2;
        pchMessageStart[2] = 0x99;
        pchMessageStart[3] = 0x88;
        vAlertPubKey = ParseHex("042d0fcbaa5d80a255cffc4c51eb34e5633f3717cb695c6f2e697f8288bfea58584cca3e60d4074ccbec9cbcd09976567cae82335dfe3142a02fe16a84378a0101");
        nDefaultPort = 21047;
        nRPCPort = 21048;
        strDataDir = "testnet";

        // Modify the testnet genesis block so the timestamp is valid for a later start.
        genesis.nTime = 1405741700;
        genesis.nNonce = 1761048;
        
     	hashGenesisBlock = genesis.GetHash();
        
        assert(hashGenesisBlock == uint256("0x00000e41ecbaa35ef91b0c2c22ed4d85fa12bbc87da2668fe17572695fb30cdf"));

        vFixedSeeds.clear();
        vSeeds.clear();
		vSeeds.push_back(CDNSSeedData("seed2.parallelcoin.info", "seed2.parallelcoin.info"));

        base58Prefixes[PUBKEY_ADDRESS] = 18;
        base58Prefixes[SCRIPT_ADDRESS] = 188;
        base58Prefixes[SECRET_KEY] = 239;

    }
    virtual Network NetworkID() const { return CChainParams::TESTNET; }
};
static CTestNetParams testNetParams;


//
// Regression test
//
class CRegTestParams : public CTestNetParams {
public:
    CRegTestParams() {
        pchMessageStart[0] = 0xf0;
        pchMessageStart[1] = 0xa1;
        pchMessageStart[2] = 0x88;
        pchMessageStart[3] = 0xcd;
        nSubsidyHalvingInterval = 150;
        bnProofOfWorkLimit[ALGO_SHA256D] = CBigNum(~uint256(0) >> 1);
        bnProofOfWorkLimit[ALGO_SCRYPT]  = CBigNum(~uint256(0) >> 1);
        genesis.nTime = 1405740500;
        genesis.nBits = 0x207fffff;
        genesis.nNonce = 1;
        hashGenesisBlock = genesis.GetHash();
        nDefaultPort = 31047;
        strDataDir = "regtest";
        
		hashGenesisBlock = genesis.GetHash();
        
        assert(hashGenesisBlock == uint256("0x69e9b79e220ea183dc2a52c825667e486bba65e2f64d237b578559ab60379181"));

        vSeeds.clear();  // Regtest mode doesn't have any DNS seeds.

        base58Prefixes[PUBKEY_ADDRESS] = 0;
        base58Prefixes[SCRIPT_ADDRESS] = 5;
        base58Prefixes[SECRET_KEY] = 128;
    }

    virtual bool RequireRPCPassword() const { return false; }
    virtual Network NetworkID() const { return CChainParams::REGTEST; }
};
static CRegTestParams regTestParams;

static CChainParams *pCurrentParams = &mainParams;

const CChainParams &Params() {
    return *pCurrentParams;
}

void SelectParams(CChainParams::Network network) {
    switch (network) {
        case CChainParams::MAIN:
            pCurrentParams = &mainParams;
            break;
        case CChainParams::TESTNET:
            pCurrentParams = &testNetParams;
            break;
        case CChainParams::REGTEST:
            pCurrentParams = &regTestParams;
            break;
        default:
            assert(false && "Unimplemented network");
            return;
    }
}

bool SelectParamsFromCommandLine() {
    bool fRegTest = GetBoolArg("-regtest", false);
    bool fTestNet = GetBoolArg("-testnet", false);

    if (fTestNet && fRegTest) {
        return false;
    }

    if (fRegTest) {
        SelectParams(CChainParams::REGTEST);
    } else if (fTestNet) {
        SelectParams(CChainParams::TESTNET);
    } else {
        SelectParams(CChainParams::MAIN);
    }
    return true;
}
