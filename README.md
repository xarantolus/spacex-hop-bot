# spacex-hop-bot
This is a [Twitter bot](https://twitter.com/wenhopbot) that informs about progress on the [SpaceX Starship](https://www.spacex.com/vehicles/starship/) by retweeting interesting tweets about it. There is a special focus on tweets from locations where Starships are built/launched, tagging tweets with a location helps the bot find them.

It reads tweets from the following sources:
* All tweets from [lists the account follows](https://twitter.com/wenhopbot/lists)
* All tweets from accounts it follows, including replies
* All tweets from [a large area around the launch and build site](https://bboxfinder.com/#25.838213,-97.321014,26.121535,-96.942673), the [SpaceX McGregor engine test site](https://mapper.acme.com/?ll=31.39966,-97.46246&z=12&t=M&marker0=31.39930%2C-97.46250%2C31.399308%20-97.462496&marker1=31.34836%2C-97.51740%2Cunnamed&marker2=31.48314%2C-97.36530%2C6.0%20km%20NE%20of%20McGregor%20TX), [Pascagoula](https://bboxfinder.com/#30.298204,-88.678894,30.457552,-88.463974), [Brownsville South Padre Island International Airport](https://bboxfinder.com/#25.891967,-97.441134,25.918835,-97.406845) and [Port/Cape Canaveral](https://mapper.acme.com/?ll=28.40952,-80.60944&z=10&t=M&marker0=28.21910%2C-80.79552%2Cunnamed&marker1=28.88617%2C-79.96262%2C79.2%20km%20ExNE%20of%20Merritt%20Island%20FL) (that are tagged with a location)

These tweets are retweeted, if:
* They contain generic keywords about Starship such as "SN11", "BN1", "Starship", "Superheavy", "raptor"
* They are from selected "trusted" users and contain info about road closures, cryogenic tests, temporary flight restrictions etc.
* They are by Elon Musk and contain anything related to Starship
* They are tagged with *exactly* the location of either [Starbase](https://twitter.com/places/1380f3b60f972001), the [Starship launch site](https://twitter.com/places/124cb6de55957000), [Starship build site](https://twitter.com/places/124bed061054f000), the [McGregor engine test site](https://twitter.com/places/07d9f642af482000), [Boca Chica Beach](https://twitter.com/places/07d9e62cfe480002) or [Boca Chica Village](https://twitter.com/places/07d9f0b85ac83003) (must have media under some criteria)

Some keywords and (mostly satire) accounts are filtered out to prevent spam. The bot tries to only retweet *real* information, which is why animations and similar are also filtered.

It also does some background tasks: 
- Watching the [SpaceX YouTube channel](https://www.youtube.com/spacex/) for livestreams. As soon as a stream about Starship goes live (or has a countdown), the bot will tweet a link.
- Checking the [Starship website](https://www.spacex.com/vehicles/starship/) from time to time to tweet if the mentioned date or Starship changed
- Checking the [Environmental Review dashboard](https://www.permits.performance.gov/permitting-project/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site) for changes and tweeting about them

You can use [this Twitter search link](https://twitter.com/search?q=from%3Awenhopbot%20-filter%3Areplies) to see these tweets.

### Contributing
If you have any suggestions for additional sources (like accounts or lists to follow) or anything else please open an issue (or write an e-mail to `x@010.one`).

If you want to edit the code, my first suggestion would be checking out the [file that defines positive and negative keywords](match/starship_keywords.go) for the matcher. The tests (run `go test ./...`) will tell you if everything still works after your changes.
