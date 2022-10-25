input="..\BashScripts\APIKey.txt"
key=`cat ${input}`
lat=$1
lon=$2
cityName=$3
countryCode=$4
limit=1
url="api.openweathermap.org/data/2.5/forecast?lat=${lat}&lon=${lon}&appid=${key}"
output="../RequestResults/WeatherData/${cityName}_${countryCode}.txt"
curl -s -f -L "${url}" -o "${output}"
return 0
