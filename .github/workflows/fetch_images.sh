#!/bin/bash
pwd
ls
for i in $(cat .github/workflows/images.txt); do
  wget "http://gwent-cards.com/wp-content/uploads/2015/06/$i" -O "static/$i"
done
