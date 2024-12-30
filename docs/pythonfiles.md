## pe.ipynb

This python file imports segno for QR code encoding, base64 and math. Ensure you do pip install segno. If you are wondering why I did not create a requirements.txt or a script.sh, I would rather people understand how this works before running it on their computer. 

A binary file is then opened as rb and saved as a variable called piece. Piece is then encoded into base64 and then chunked out in 2000kb pieces. Finally they are saved to the img folder. 

The second chunk is for me to verify hashes and to confirm that if the binary is put back together the hash is the same. 

### Quick lesson on hashing
Hashing is 3 main purposes. 
1) A hash has to be unique. That means that no matter what I put in, there is a unique fingerprint associated with it that doesn't resemble a similar value. I.e the hash of a has to look completely different to the hash of aa.
2) A hash has to go one way. This means it should be impossible to reverse a hash to it's original value. Companies like to store hashes of passwords instead of the passwords themselves because if a hacker breaches that list, they will not be able to reverse it back. Sure there's salting and rainbow tables that I could talk about but I don't want to type for longer than I need to.
3) A hash has to have the same length of every other hash. For example, every sha256 hash will have a length of 64 characters. 


## Caveats
finalinfectedwithhardcoded pid is used here. To get the right pid make sure you run ./binaries/pid and replace that variable. I did it like this because I couldn't be bothered to figure out how to use command arguments while doing process injection. 



