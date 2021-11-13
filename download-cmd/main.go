package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/willena/S3Replicator/manifest"
)

func main() {
	log.Println("OK !")

	//Golang.zip
	//d3d6fc57adde9749107655a998c25a19e9b687455fc688b70788883f18c1f94c,Goland/Goland.zip,4,1,'0967F76D0B510E6E281DDDB6281AB352C14705B931601888E2030A35F76F98B1',134283264
	//666573fe676c220a234064649731762a7845f99276045576ae7796d8a8d09de1,Goland/Goland.zip,4,2,'22A870E279BB4960064A12E1F22593B195E0B884F9E24325B7DB3CE6AF5B275C',134283264
	//5e5893ebc2d8685d571e8333b8e5d6bc7cedb451ee75932ebef186436e613256,Goland/Goland.zip,4,3,'A5BF22BE9DCE5D4C342DCCF6E83E7E3F12D7510DE0F3DDF25138B2782EE73059',134283264
	//8e7f5994e4212d7b65909610a7bfa78438c7f9c3c9b79940ace6849a0ed2dd1b,Goland/Goland.zip,4,4,'76B4EE7FA7BFCFBA00E626D4EC0B63DF14C651D9C7541F7C832579213215DFCA',128517520

	//c5c7cce46c20885c276ccf73f9b9d4eacb15d05a75e09c09af9368c1dd150bfb,Goland/App/jbr/include/jawt.h,1,1,'CA762FBD3DE6818E8C683B2A7D95570273C5A80EEA1641D1BC9A105B27810B22',12490

	headerFileItem := manifest.Get().GetObject("c5c7cce46c20885c276ccf73f9b9d4eacb15d05a75e09c09af9368c1dd150bfb")
	log.Debug("Item /jawt.h from obj", headerFileItem)

	golangZipItem3 := manifest.Get().GetObject("5e5893ebc2d8685d571e8333b8e5d6bc7cedb451ee75932ebef186436e613256")
	log.Debug("Item Goland.zip _3 : ", golangZipItem3)

	allPartsFromAnyObjectId_3 := manifest.Get().GetPartsFromObjectId("5e5893ebc2d8685d571e8333b8e5d6bc7cedb451ee75932ebef186436e613256")
	log.Debug("All Goland.zip part from p3 : ", allPartsFromAnyObjectId_3)

	allPartsFromFileName := manifest.Get().GetPartsFromFileName("Goland/Goland.zip")
	log.Debug("All Goland.zip parts : ", allPartsFromFileName)


}
