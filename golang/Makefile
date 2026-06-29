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

# ASN.1 source file - uses ASNFILES from root Makefile (exported)
ASN_SOURCE = $(addprefix ../, $(ASNFILES))

# rule to convert the asn1 to Go code
# Falls back to existing manual structs if asn1go cannot handle the new spec
generate:
	@echo "Attempting asn1go code generation from $(ASN_SOURCE)..."
	@if cd asn1go/cmd/asn1go && go run main.go ../../../$(notdir $(ASN_SOURCE)) ../../../ipmtscd/IPMTSCD.go 2>/dev/null; then \
		echo "asn1go generation succeeded."; \
	else \
		echo "WARNING: asn1go failed (likely due to AUTOMATIC TAGS or extensibility markers)."; \
		echo "Using existing manually maintained structs in ipmtscd/IPMTSCD.go"; \
	fi
	@echo "Verifying Go build..."
	go build
	@echo "Go build verification passed."

# rule to build the Go application
build:
	go build

# rule to run the decode and then encode
run:
	./testapp -m recode

.PHONY: generate build run
