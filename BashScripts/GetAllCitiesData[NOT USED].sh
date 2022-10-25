input="..\BashScripts\Cities.txt"

while read -r line
do
  echo "$line"
  vars=( $line )
  ./GetCityData.sh ${vars[0]} ${vars[1]} 
done < "$input"

