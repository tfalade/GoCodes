# GoCodes

A GET call to the endpoint <counties> like below will return all the elections results for all the fips( Federal Information Processing Standard) as specified in the csv files.

GET http://your-service-here/counties Returns all counties with election results in the following format:

{ "counties": [ { "name": "Adair", "fips": "19001", "elections": [ { "party": "Democratic", "results": [ {"candidate": "Hillary Clinton", "votes": 113}, {"candidate": "Bernie Sanders", "votes": 86}, {"candidate": "Martin O'Malley","votes": 0} ] } { "party": "Republican", "results": [ {"candidate": "Ted Cruz","votes": 104}, {"candidate": "Donald Trump","votes": 104} ... ] } ] }, { "name": "Adams", "fips": "19003", "elections": [ ... ] } ] }

You can also choose a particular fips where you want to see the election result. For example fip 19001, make the GET call to the endpoint as stated below and you wull see the results.

GET http://your-service-here/counties/19001  Returns the results for the county with the specified FIPS code.

{ "name": "Adair", "fips": "19001", "elections": [ { "party": "Democratic", "results": [ {"candidate": "Hillary Clinton", "votes": 113}, {"candidate": "Bernie Sanders","votes": 86}, {"candidate": "Martin O'Malley","votes": 0} ] } { "party": "Republican", "results": [ {"candidate": "Ted Cruz","votes": 104}, {"candidate": "Donald Trump","votes": 104} ... ] } ] }
