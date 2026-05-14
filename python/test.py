# MIT License
#
# Copyright (c) 2023 Yunex Limited
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.


import asn1tools
import base64
import hashlib
import uuid
import time
import json
import sys

print("Loading Compiler...")

jsonFile= sys.argv[1]
asnlist = sys.argv[2:]

print(asnlist)

def ordered(obj):
    if isinstance(obj, dict):
        return sorted((k, ordered(v)) for k, v in obj.items())
    if isinstance(obj, list):
        return [ordered(x) for x in obj]
    if isinstance(obj, tuple):
        return tuple(ordered(x) for x in obj)
    else:
        return obj

# load in the reference file
jerEncoded = asn1tools.compile_files(asnlist, 'jer')

# load in the reference file
file = open(jsonFile, "r") 

txtdata = file.read()

#print("Src json: "+str(txtdata))

# decode to internal format
decodedMessage_reference = jerEncoded.decode("IpmstscdData", bytearray(txtdata, encoding='utf8'))

# and re-encode 
reEncoded = jerEncoded.encode("IpmstscdData", decodedMessage_reference).decode("utf-8") 
                              
#print("reEncoded json: "+str(reEncoded))

srcJson=json.loads(txtdata)
reEncodedJson=json.loads(reEncoded)

print("Src json: "+str(srcJson))
print("reE json: "+str(reEncodedJson))

if srcJson != reEncodedJson:
    raise Exception("Reference JSON has felds that cannot be parsed by standard") 

#--------------------------------
# and define the tests
tests = [
            {
                "srcfile":"../c-detector.ber", 
                "binary": True,
                "decoder":"ber"
            },
            {
                "srcfile":"../golang-detector.ber", 
                "binary": True,
                "decoder":"ber"
            }
        ]

print("Reference: "+str(decodedMessage_reference))

# for each test....
for curTest in tests:
    
    print("##################################################")
    print("Performing test on "+curTest["srcfile"]+" using "+curTest["decoder"])
    
    decoder  = asn1tools.compile_files(asnlist, curTest["decoder"])
    
    
    openmode="r"
    if curTest["binary"]:
        openmode=openmode+"b"
    
    file = open(curTest["srcfile"], openmode) 
    
    txtdata = file.read()
    
    if curTest["binary"]==False:
        txtdata=bytearray(txtdata, encoding='utf8') 
    
    # decode to internal format
    decodedMessage_c = decoder.decode("IpmstscdData", txtdata)
    
    print(str(decodedMessage_c))
    
    # and test them
    if ordered(decodedMessage_c) != ordered(decodedMessage_reference):
        raise Exception("Test does not pass") 


