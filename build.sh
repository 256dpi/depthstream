#!/usr/bin/env bash

rm -f depthstream
rm -f depthstream.tar.gz

gopm build

mv output depthstream

tar -czf depthstream.tar.gz depthstream
