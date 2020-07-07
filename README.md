# terraform-module-readme-gen

## Usage

```bash
terraform-module-readme-gen -readme README.md -variables var.tf -outputs out.tf 
```

then you can check the `README.md` file content.

## Note

1.This tool only generates the boring part which describes variables and output fields, you need to modify the left parts according your specific module information.

2.The output of the tool **IS NOT 100% CORRECT** but could leave the most doc editing work out, please **RECHECK** the content before you release it. 

## Thanks
Thanks to [hcl2json](https://github.com/tmccombs/hcl2json) for leading us a way for the conversion.