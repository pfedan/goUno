# goUno
A Go implementation of a popular card game: UNO

**Contributors are welcome**

## Motivation
This project was started with a variety of intentions:
- **Shuffling**: I wanted to analyze the [Riffle]([http://](https://en.wikipedia.org/wiki/Shuffling#Riffle)) technique for its randomness and its limits.
- **Strategies**: Is it advantageus to play aggressively / stay with a certain color / switch colors if possible? In order to analyse new strategies, one can implement more cases in `scoreCandidates(...)` in [player.go](goUno/player.go). 
- **Entertainment**: Just play UNO against a computer or together with other human players. 
- **Web readiness**: It will be rather straight forward to integrate the game engine as a server towards a web frontend that provides the game to users that wish for a more convenient UI than CLI.

## Usage
Clone the repository and build the binary:

```bash
git clone https://github.com/pfedan/goUno.git
cd goUno
go build
```

Run the game (using different parameters):
```bash
./goUno.exe # Default: Two players "A" and "B", non-human, 1 round
./goUno.exe -p Adam,Berta,Charles # Three players with custom names
./goUno.exe -p Human,Computer1,Computer2 -h 1,0,0 # Three players, one human player defined with the -h flag list
```

Analyze a large number of rounds (example):
 ```bash
# Command
./goUno.exe -m -r 100000 -p A,B,C,D # Muted game log (-m) and 100k rounds with four players: A, B, C and D

# Output:
Wins per Player: map[A:27053 B:24346 C:24122 D:24479]
Total points per Player: map[A:3150105 B:2749662 C:2714910 D:2762753]
Count of turns per game: map[15:1 16:7 17:22 18:99 19:258 20:567 21:942 22:1391 23:1734 24:1918 25:1900 26:1937 27:1922 28:1997 29:2029 30:1986 31:2059 32:2076 33:2110 34:2090 35:2047 36:2018 37:2049 38:2088 39:2020 40:2041 41:2142 42:2104 43:2013 44:1944 45:1923 46:1951 47:1880 48:1778 49:1664 50:1769 51:1670 52:1647 53:1563 54:1570 55:1511 56:1506 57:1410 58:1371 59:1281 60:1204 61:1227 62:1203 63:1161 64:1102 65:1032 66:992 67:959 68:860 69:893 70:826 71:849 72:744 73:731 74:701 75:649 76:635 77:622 78:567 79:536 80:559 81:514 82:467 83:440 84:447 85:395 86:390 87:381 88:361 89:334 90:343 91:329 92:287 93:282 94:265 95:231 96:253 97:220 98:226 99:218 100:172 101:197 102:178 103:170 104:156 105:147 106:145 107:122 108:135 109:118 110:99 111:120 112:94 113:103 114:87 115:68 116:72 117:91 118:70 119:72 120:59 121:61 122:58 123:68 124:58 125:47 126:54 127:38 128:40 129:36 130:31 131:28 132:35 133:25 134:27 135:25 136:28 137:27 138:27 139:16 140:16 141:17 142:13 143:15 144:11 145:9 146:14 147:15 148:15 149:20 150:14 151:4 152:11 153:5 154:12 155:8 156:5 157:6 158:5 159:6 160:8 161:7 162:1 163:5 164:7 165:4 166:6 167:3 168:7 169:3 170:4 171:7 172:8 173:2 174:5 175:2 176:3 177:2 178:3 179:2 181:1 182:1 183:2 184:2 185:5 186:3 187:2 188:1 189:1 190:2 191:1 192:3 194:3 195:1 196:3 198:1 199:1 200:2 201:1 202:1 203:1 206:4 207:1 209:1 212:1 228:1]
Count of points per game: map[15:1 16:7 17:22 18:99 19:258 20:567 21:942 22:1391 23:1734 24:1918 25:1900 26:1937 27:1922 28:1997 29:2029 30:1986 31:2059 32:2076 33:2110 34:2090 35:2047 36:2018 37:2049 38:2088 39:2020 40:2041 41:2142 42:2104 43:2013 44:1944 45:1923 46:1951 47:1880 48:1778 49:1664 50:1769 51:1670 52:1647 53:1563 54:1570 55:1511 56:1506 57:1410 58:1371 59:1281 60:1204 61:1227 62:1203 63:1161 64:1102 65:1032 66:992 67:959 68:860 69:893 70:826 71:849 72:744 73:731 74:701 75:649 76:635 77:622 78:567 79:536 80:559 81:514 82:467 83:440 84:447 85:395 86:390 87:381 88:361 89:334 90:343 91:329 92:287 93:282 94:265 95:231 96:253 97:220 98:226 99:218 100:172 101:197 102:178 103:170 104:156 105:147 106:145 107:122 108:135 109:118 110:99 111:120 112:94 113:103 114:87 115:68 116:72 117:91 118:70 119:72 120:59 121:61 122:58 123:68 124:58 125:47 126:54 127:38 128:40 129:36 130:31 131:28 132:35 133:25 134:27 135:25 136:28 137:27 138:27 139:16 140:16 141:17 142:13 143:15 144:11 145:9 146:14 147:15 148:15 149:20 150:14 151:4 152:11 153:5 154:12 155:8 156:5 157:6 158:5 159:6 160:8 161:7 162:1 163:5 164:7 165:4 166:6 167:3 168:7 169:3 170:4 171:7 172:8 173:2 174:5 175:2 176:3 177:2 178:3 179:2 181:1 182:1 183:2 184:2 185:5 186:3 187:2 188:1 189:1 190:2 191:1 192:3 194:3 195:1 196:3 198:1 199:1 200:2 201:1 202:1 203:1 206:4 207:1 209:1 212:1 228:1]
```

Plotting the above counts of turns per game as relative and cumulative frequencies will result in the following graphs:

![Relative and cumulative frequency of turns](doc/example_4players_100k_rounds.png)

![Relative and cumulative frequency of points](doc/example_4players_100k_points.png)

For instance you can see, that if everyone plays aggressively, the winner gets more points at the end of a round:
![Random vs. aggressive strategy](doc/example_4players_100k_points_rand_vs_aggressive.png)