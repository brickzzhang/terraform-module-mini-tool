# terraform-module-readme-gen

## Usage

1.Download the binary file from [latest workflow](https://github.com/brickzzhang/terraform-module-readme-gen/actions/).

2.Run the binary file as below:
```bash
terraform-module-readme-gen_darwin_amd64 -readme DEMO-README.md -variables var.tf -outputs out.tf 
```

then you can copy the content of `DEMO-README.md` file.

## Note

1.This tool only generates the boring part which describes variables and output fields, you need to modify the left parts according your specific module information.

2.The output of the tool **IS NOT 100% CORRECT** but could leave the most doc editing work out, please **RECHECK** the content before you release it. 

3.Specifying `-readme` to your **actual module README.md file** will overwrite it, it's recommended to specify it to other file, such as `DEMO-README.md`.  

## APPRECIATION
Thanks to [hcl2json](https://github.com/tmccombs/hcl2json) for leading us a way for the conversion.