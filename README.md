# terraform-module-readme-gen

## Usage

```bash
terraform-module-readme-gen -readme DEMO-README.md -variables var.tf -outputs out.tf 
```

then you can copy the content of `DEMO-README.md` file.

## Note

1.This tool only generates the boring part which describes variables and output fields, you need to modify the left parts according your specific module information.

2.The output of the tool **IS NOT 100% CORRECT** but could leave the most doc editing work out, please **RECHECK** the content before you release it. 

3.Specifying `-readme` to your **actual module README.md file** will overwrite it, it's recommended to specify it to other file, such as `DEMO-README.md`.  

## APPRECIATION
Thanks to [hcl2json](https://github.com/tmccombs/hcl2json) for leading us a way for the conversion.