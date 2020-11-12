package testsample

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var PublicKeys = []hexutil.Bytes{
	hexutil.MustDecode("0x0af4b9a4e9e2b5a4d4d0a6d2eb9af19abd9d8c5f009b50f3d15faa7e7064f69f"),
	hexutil.MustDecode("0x5db26f64f4453a0f3182fe5f07a9f3ac862aa2ae679ca424a66d999121ed5ca6"),
	hexutil.MustDecode("0x8f21da423cceb451c280dd50ae017ac89ae508779dbaec52f6b6a9e98bfaa4a3"),
	hexutil.MustDecode("0x4554776709196eff2a48fdc235d5192385fab31aa66413d3757d3c3b4ceb9482"),
	hexutil.MustDecode("0xb974e2f25ed9bc54c6d5cbefcd134d3e8c8ac24fd18c1c424eb650993be2b12b"),
	hexutil.MustDecode("0xf877a4103cd6192a159602b4771944edd705c074227285119b6b65485419c20c"),
	hexutil.MustDecode("0xedb615e8cd65fd90c95036613890642c11e56f031fbbac9575d9559a0cc0b588"),
	hexutil.MustDecode("0x71515c496daccd10e991f508952711a76176854794f90c4473364d177d72da97"),
	hexutil.MustDecode("0xd63c298ee09d12488796602d4591854e23424f28a8a73dd7df1880721311a88b"),
	hexutil.MustDecode("0xdedbbb5964853d19a99875b2cdd4640c42b3cac5d4469c91a8f53c2f5b3c45af"),
}

var PrivateKeys = []hexutil.Bytes{
	hexutil.MustDecode("0x05befa1dc5beb8aa74c348966f5254702bc0a9613e519eb3ef2fe8c444f40d33"),
	hexutil.MustDecode("0x01dc5faec310a1c6d1e939a78f815de4bfd13c7fbf5f88055eab1b19df44fdc4"),
	hexutil.MustDecode("0x02de0338d90531e485fbf7f78cb16c2a7d607ad2c27eccf18eb13a9255775c2e"),
	hexutil.MustDecode("0x0233da04457bf4c0a116675a272ecdbd3089ff7d212a2416762b0b2589cc2ce3"),
	hexutil.MustDecode("0x0078ec4f65ca8b49748c3175213ff52c675cbb8aad9a1ad0275b6f83e352d4a0"),
	hexutil.MustDecode("0x00e5ccfba6f1d3ab314d6d3217b9048effd8365f7c936c3a93ea2600d0325356"),
	hexutil.MustDecode("0x01ee0c930ef3f8efd69546043b06ee9697bf86c13bbf7b6324290cf8741503d2"),
	hexutil.MustDecode("0x00e6f6250d858dc21d8e7bd0a2f9b6027ad1f47abdc06e8f0de9c9fdecfc110a"),
	hexutil.MustDecode("0x02ed199548ec3ce6e34474f8aa7654c345287ea87af36ba4d9bcebbb4621c506"),
	hexutil.MustDecode("0x01cf655a41722cbbe3c8921dcb0daaf22bb1ae3bc008d5f342a9a91a2618238b"),
}

var Accounts = []common.Address{
	common.HexToAddress("0x4190b3B9A0fcA9B1A2f4C77B22A2b32F40224177"),
	common.HexToAddress("0x41906015d064Ad593Fba0c6deC0c714bCB18f269"),
	common.HexToAddress("0x419021B62197c40b081FABc0E7E36910B11F27dD"),
	common.HexToAddress("0x419086388E5073F581403730a9b85bB26a944Ab2"),
	common.HexToAddress("0x4190D2209718a6c8472cB03706A57f8e1f6b9589"),
	common.HexToAddress("0x4190e044505b27Db6028865B1f48A4398a68d887"),
	common.HexToAddress("0x4190Ea9bad72fB2a38DD96Bc1D41241Bf7675382"),
	common.HexToAddress("0x4190Ebc1C0E79DdD1DDFec0e9521663Ff9a15c50"),
	common.HexToAddress("0x41901e9e37271b587dbC978e9Fb369Ede85bf021"),
	common.HexToAddress("0x41908d6b087C97337a7b5218cbA9c69d391C16B2"),
}
