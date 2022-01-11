package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Data struct {
	Builds []struct {
		Name          string      `json:"name"`
		BuilderType   string      `json:"builder_type"`
		BuildTime     int         `json:"build_time"`
		Files         interface{} `json:"files"`
		ArtifactID    string      `json:"artifact_id"`
		PackerRunUUID string      `json:"packer_run_uuid"`
		CustomData    struct {
			MyCustomData string `json:"my_custom_data"`
		} `json:"custom_data"`
	} `json:"builds"`
	LastRunUUID string `json:"last_run_uuid"`
}

func executePacker(packerFile string) {
	//exec packer
	cmd := exec.Command("packer", "build", packerFile)
	println("\n-------Packer build-------")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
}

func writeAmiIdInTFVars(manifest string, tfvarFile string) {
	println("\n-------opening manifest.json generated by packer-------\n")
	//opening manifest.json generated by packer
	jsonFile, err := os.Open(manifest)

	if err != nil {
		fmt.Printf("failed to open json file: %s, error: %v", manifest, err)
		return
	}
	defer jsonFile.Close()
	println("\n-------reading and validating manifest.json-------\n")
	//reading and validating manifest.json
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return
	}
	println("\n-------unmarshalling manifest.json-------\n")
	//unmarshalling manifest.json
	data := Data{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}
	println("\n-------writing the AMI in varibles.tf-------\n")
	// Print
	for _, item := range data.Builds {
		ami := strings.Split(item.ArtifactID, ":")
		fmt.Printf("artifact_id: %s \n", ami[1])
		//opening manifest.json generated by packer
		content := "\n//AMI exported by packer\nvariable \"ami\" {\n  default     = \"" + ami[1] + "\"\n  description = \"The latest AMI.\"\n}"
		var_tf, err := os.OpenFile(tfvarFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println(err)
		}
		defer var_tf.Close()

		_, err2 := var_tf.WriteString(content)

		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

func executeTerraform(path string, action string) {

	init := exec.Command("terraform", "-chdir=./terraform", "init")
	println("\n-------Terraform initializing-------")
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	init_err := init.Run()

	if init_err != nil {
		log.Fatal(init_err)
	}
	init.Wait()

	apply := exec.Command("terraform", "-chdir=./terraform", action, "-auto-approve")

	println("\n-------Terraform apply-------")
	apply.Stdout = os.Stdout
	apply.Stderr = os.Stderr
	apply_err := apply.Run()

	if apply_err != nil {
		log.Fatal(apply_err)
	}
	apply.Wait()
}

func main() {

	executePacker("./packer/aws-ubuntu.pkr.hcl")
	writeAmiIdInTFVars("./packer/manifest.json", "./terraform/variables.tf")
	executeTerraform("./terraform", "apply")

}
