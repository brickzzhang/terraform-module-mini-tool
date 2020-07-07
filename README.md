# terraform-module-readme-gen

## Usage

This script is used to generate simple readme demo of tencentcloud terraform modules.

1.Download script under binary directory.

2.Run the script, then input your path of variables.tf and outputs.tf as it asked, such as below:

`/Users/brick/workspace/terraform/terraform-tencentcloud-modules/terraform-tencentcloud-clb/`

3.Check DEMO-README.md file under current directory and copy the content to your target file.

## Note

1.This script only generates the boring part which describes variables and output fields, you need to modify the left parts according your specific module information.

2.The output of the script is NOT 100% CORRECT but could leave the most doc editing work out, please RECHECK the content before you release it. 

## Thanks
Thanks to [hcl2json](https://github.com/tmccombs/hcl2json) for leading us a way for the conversion.