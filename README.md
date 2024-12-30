# Documentation

## Introduction
the purpose of this project is a novel approach to malware dropping. While the implementation will not generally work in the real world, it aims to show how it can be possible to initially drop malware via QR codes. 


## Project structure

binaries -- Contains all the binaries that I used to develop the project

code -- Contains c code that I used as inspiration 

docs -- documentation in depth.

elf_to_shellcode_master -- Is a python project that I took from github to try turn binaries into shellcode. Didn't really work.

gofiles -- contains all the go scripts created to get this to work.

gozxing-master -- Library used to decode qr codes

images -- holds some dummy images

img -- holds images of the binary that does process injection

qr_tx-master -- I forgot why this is included

pe.ipynb -- python file that breaks a binary into base64 chunks and converts them into QR codes. 

## open source libraries used in the  project
Gozxing -- QR code decoder
Segno -- QR code encoder in python which allowed me to use v40 
Elf-to_shellcode -- This was used to try turn a binary into "shellcode" but it didn't really work but I'm leaving it in here in case anyone wants to try get it to work.

