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

# Python interpreter — exported by the root Makefile. Default to python3.11
# if invoked directly without the parent Makefile.
PYTHON ?= python3.11

../python-detector.ber: ../${SRC_JSONFILE}
	@${PYTHON} -m venv venv
	@. venv/bin/activate && pip install -r requirements.txt
	@. venv/bin/activate && python det_encoder.py ../${SRC_JSONFILE} $(addprefix ../, ${ASNFILES}) 
	

# rule to test the different files generated
test:
	@${PYTHON} -m venv venv
	@. venv/bin/activate && pip install -r requirements.txt
	@. venv/bin/activate && python test.py ../${SRC_JSONFILE} $(addprefix ../, ${ASNFILES}) 

