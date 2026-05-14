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


ASNFILES=Annex_A_10711_v3.0.asn

SRC_JSONFILE=detection_v3.json

# Python interpreter used by sub-Makefiles. Override on the command line if needed
# (e.g. `make PYTHON=python3.12 run_python`).
PYTHON=python3.11


export ASNFILES
export SRC_JSONFILE
export PYTHON

all: 

# rule to make the docker image
docker:
	docker build -t topas-2545 . 


# Rule to convert the ASN into the language specific files
generate_all:
	# Generate the Python version - Not required
	# Generate the C version
	${MAKE} generate_c
	# Generate the Golang version 
	${MAKE} generate_go

# Rule to convert the ASN1 into the language files
build_all:
	# Build the Python version - Not required
	# Build the C version
	${MAKE} build_c
	# Build the Golang version 
	${MAKE} build_go

# rule to run the differnet Languages
run_all:
	${MAKE} run_python
	${MAKE} run_c
	${MAKE} run_go

test_all:
	${MAKE} test_python


clean_all:
	${MAKE} -C c clean
	rm -f c-detector.* python-detector.* golang-detector.*

# Rule to generate the C code from ASN1
generate_c:
	@${MAKE} -C c generate

# Rule to generate the GO code from ASN1
generate_go:
	@${MAKE} -C golang generate

# Rule to generate and build the C version
build_c:
	@${MAKE} -C c build	

# Rule to build the GOlang version	
build_go:
	@${MAKE} -C golang build	

# Run the Python version
run_python:
	@${MAKE} -C python

# Run the C version
run_c:
	# check that we can parse the ber formated data into XML
	@c/converter-sample -iber -oxer python-detector.ber > c-detector.xer
	# we then recode it as DER (DER is a subset of BER so are compatable)
	@c/converter-sample -ixer -oder c-detector.xer > c-detector.ber

run_go:
	@${MAKE} -C golang run
	

# rule to use the python version to test the c and golang versions
test_python:
	@${MAKE} -C python test



	
.PHONY: all run_python run_c run_go


