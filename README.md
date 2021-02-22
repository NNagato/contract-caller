# contract-caller

### Setup

**Step 1: Run go app**<br/>
Go to cmd directory and run these following commands:
+ ```go build -v .```
+ run ```./cmd``` with your configuration flags

**Step 2: Run react app**<br/>
Go to html/app directory and run these following commands:
+ ```npm install```
+ ```npm start```

**Step 3: Enjoy the app on browser**<br/>
Open  your browser and enter this url ```http://localhost:3000```

### Supported Data Types

- ```address``` e.g input: ```0x1e7a39bc29e07fc214646c3574aba8a2dbefdad1```
- ```address[]``` e.g input: ```0x1e7a39bc29e07fc214646c3574aba8a2dbefdad1, 0xa090e606e30bd747d4e6245a1517ebe430f0057e, ...```
- ```uint256``` e.g input: ```123456```
- ```uint256[]``` e.g input: ```123456, 234567```
- ```bool``` e.g input: ```true```
- ```bool[]``` e.g input ```true, false, ...```
- ```bytes```, ```int8```, ```bytes32``` input should be in hex format, e.g input ```0x123```
