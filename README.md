# spacex-hop-bot
This is a [Twitter bot](https://twitter.com/wenhopbot) that informs about progress on the [SpaceX Starship](https://www.spacex.com/vehicles/starship/) by retweeting interesting tweets about it.

It reads tweets from the following sources:
* All tweets from [lists the account follows](https://twitter.com/wenhopbot/lists)
* All tweets from accounts it follows, including replies
* All tweets from [around the launch area](https://mapper.acme.com/?ll=26.00002,-97.07932&z=10&t=M&marker0=25.98750%2C-97.18639%2CSpaceX%20South%20Texas%20launch%20site&marker1=26.39190%2C-96.71811%2C26.3919%20-96.7181&marker2=25.52629%2C-97.43501%2C25.5263%20-97.4350) (that are tagged with a location)

These tweets are retweeted, if:
* They contain generic keywords about Starship such as "SN11", "BN1", "Starship", "Superheavy", "raptor"
* They are from selected "trusted" users and contain info about road closures, cryogenic tests, temporary flight restrictions etc.
* They are by Elon Musk and contain anything related to Starship

Some keywords are filtered out to prevent spam.

It also checks the [SpaceX YouTube channel](https://www.youtube.com/spacex/) for livestreams. As soon as a stream about Starship goes live, the bot will tweet a link.

If you have any suggestions for additional tweet sources or anything else please open an issue. 
If you have suggestions for accounts or lists to follow, you can write the bot a DM on Twitter or open an issue here.
