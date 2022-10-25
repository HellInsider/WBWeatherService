input="..\BashScripts\APIKey.txt"
key=`cat ${input}`
cityName=$1
countryCode=$2
limit=1
url="http://api.openweathermap.org/geo/1.0/direct?q=${cityName},'$'${countryCode}&limit=${limit}&appid=${key}"
output="../RequestResults/CityData/${cityName}_${countryCode}.txt"
echo $(curl -s -f -L "${url}" -o "${output}")
return 0
