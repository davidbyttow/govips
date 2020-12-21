ksdiff "$1" "$2"
read -p "Does this image match? $1 (y/n):" -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    mv "$2" "$1"
fi