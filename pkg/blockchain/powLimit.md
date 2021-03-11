# tracing powLimit

pkg/chain/validate.go checkProofOfWork

- pkg/chain/validate.go checkBlockHeaderSanity
  - pkg/chain/validate.go checkBlockSanity
    - pkg/chain/process.go ProcessBlock
    - pkg/chain/validate.go CheckBlockSanity
    - pkg/chain/validate.go CheckConnectBlockTemplate *

- pkg/chain/validate.go CheckProofOfWork
  - cmd/node/getwork.go handleGetWorkSubmission
