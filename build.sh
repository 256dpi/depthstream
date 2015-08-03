#!/usr/bin/env bash

rm depthstream.tar.gz

gopm bin github.com/256dpi/depthstream
tar -cvzf depthstream.tar.gz depthstream
