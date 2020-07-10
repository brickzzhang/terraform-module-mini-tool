# terraform-module-readme-gen

## terraform-module-mini-tool

A mini tool which support below functions:

 - Module initiation
 Create necessary files for terraform module. 
 
 - Readme generation
 Generate readme file automatically according to variables.tf and outputs.tf.

## Usage

1.Download the binary file from [release assets](https://github.com/brickzzhang/terraform-module-readme-gen/releases).

2.Run the binary file and follow the hint:

```markdown
|---------------------------------------------------|
|------ Please select which function you want ------|
|------ [0] module init ----------------------------|
|------ [1] README generate ------------------------|
|------ [q] quit -----------------------------------|
|---------------------------------------------------|
```

then you can copy the content of `DEMO-README.md` file under your specific dir.

## Note

1.This tool only generates the boring part which describes variables and output fields, you need to modify the left parts according your specific module information.

2.The output of the tool **IS NOT 100% CORRECT** but could leave the most doc editing work out, please **RECHECK** the content before you release it. 

3.Specifying `-readme` to your **actual module README.md file** will overwrite it, it's recommended to specify it to other file, such as `DEMO-README.md`.  

## APPRECIATION
Thanks to [hcl2json](https://github.com/tmccombs/hcl2json) for leading us a way for the conversion.