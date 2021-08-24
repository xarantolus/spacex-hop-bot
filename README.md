# spacex-hop-bot
This is a [Twitter bot](https://twitter.com/wenhopbot) that informs about progress on the [SpaceX Starship](https://www.spacex.com/vehicles/starship/) by retweeting interesting tweets about it.

It reads tweets from the following sources:
* All tweets from [lists the account follows](https://twitter.com/wenhopbot/lists)
* All tweets from accounts it follows, including replies
* All tweets from [a large area around the launch and build site](https://mapper.acme.com/?ll=26.00002,-97.07932&z=10&t=M&marker0=25.98750%2C-97.18639%2CSpaceX%20South%20Texas%20launch%20site&marker1=26.39190%2C-96.71811%2C26.3919%20-96.7181&marker2=25.52629%2C-97.43501%2C25.5263%20-97.4350) and the [SpaceX McGregor engine test site](https://mapper.acme.com/?ll=31.39966,-97.46246&z=12&t=M&marker0=31.39930%2C-97.46250%2C31.399308%20-97.462496&marker1=31.34836%2C-97.51740%2Cunnamed&marker2=31.48314%2C-97.36530%2C6.0%20km%20NE%20of%20McGregor%20TX) (that are tagged with a location)

These tweets are retweeted, if:
* They contain generic keywords about Starship such as "SN11", "BN1", "Starship", "Superheavy", "raptor"
* They are from selected "trusted" users and contain info about road closures, cryogenic tests, temporary flight restrictions etc.
* They are by Elon Musk and contain anything related to Starship
* They are tagged with *exactly* the location of either the [build site](https://twitter.com/places/124bed061054f000), [launch site](https://twitter.com/places/124cb6de55957000) or [Starbase](https://twitter.com/places/1380f3b60f972001)

Some keywords and (mostly satire) accounts are filtered out to prevent spam. Tweets with too many selfies are also filtered out.

The bot tries to only retweet *real* information, which is why animations and similar are also filtered.

It also does some background tasks: 
- Watching the [SpaceX YouTube channel](https://www.youtube.com/spacex/) for livestreams. As soon as a stream about Starship goes live (or has a countdown), the bot will tweet a link.
- Checking the [Starship website](https://www.spacex.com/vehicles/starship/) from time to time to tweet if the mentioned date or Starship changed

If you have any suggestions for additional sources or anything else please open an issue. 
If you have suggestions for accounts or lists to follow, you can write the bot a DM on Twitter or open an issue here.
