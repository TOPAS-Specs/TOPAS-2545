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

xml_decoder = asn1tools.compile_files(asnlist, codec='xer')
jerEncoded = asn1tools.compile_files(asnlist, 'jer')
uperEncoded = asn1tools.compile_files(asnlist, 'uper')
derEncoder = asn1tools.compile_files(asnlist, 'der')
berEncoder  = asn1tools.compile_files(asnlist, 'ber')

# get the RAW XML data
file = open(jsonFile, "r") 

txtdata = file.read()

# decode to internal format
decodedMessage = jerEncoded.decode("IpmstscdData", bytearray(txtdata, encoding='utf8'))

print(decodedMessage)

res = berEncoder.encode("IpmstscdData", decodedMessage)

print(res)

with open("../python-detector.ber", "wb") as newfile:
    newfile.write(bytearray(res))



res = xml_decoder.encode("IpmstscdData", decodedMessage)

print(res)

with open("../python-detector.xer", "wb") as newfile:
    newfile.write(bytearray(res))
    

res = derEncoder.encode("IpmstscdData", decodedMessage)

print(res)

with open("../python-detector.der", "wb") as newfile:
    newfile.write(bytearray(res))
    
    
    