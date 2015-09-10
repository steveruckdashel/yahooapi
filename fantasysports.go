//Fantasy Sports API
//
// Documentation from https://developer.yahoo.com/fantasysports/guide
package yahooapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "encoding/json"
	"encoding/xml"
	"os"
	"golang.org/x/oauth2"
)

// `json:"myName,omitempty"`
// `json:"-"`
// `json:",string"` // float/int stored in json string
// `xml:"Group>Value"`
// `xml:"where,attr"`  // pull from the attribute "where"

// Introduction to Fantasy Sports API
//
// The Fantasy Sports APIs provide URIs used to access fantasy sports data.
// Currently the APIs support retrieval of Fantasy Football, Baseball,
// Basketball, and Hockey data including game, league, team, and player
// information. The APIs are based on a RESTful model. Therefore, resources like
// game, league, team, player etc. and collections like games, leagues, teams,
// players form the building blocks for these APIs. Each resource is
// identified by a resource ID, and a collection is identified by its scope,
// specified in the URI.
//
// Historically, Yahoo! has provided two full draft and trade style fantasy
// football and baseball games – a free version, and a plus version (which
// contains more features and content). With the 2010 seasons, the Free and
// Plus versions of Football and Baseball have merged. Each game is comprised
// of many “leagues”, which typically contain 8-12 teams, which are managed by
// one or more users. At the beginning of a league’s season, professional
// athletes (“players”) are uniquely assigned or chosen through a draft to each
// team. The players that are not chosen or assigned are available to be
// acquired via a free-agent or waiver wire process (a “transaction”). These
// teams compete against each other based on statistics from real-world
// competitions based on categories like touchdowns, yards gained, batting
// average and ERA. Many fantasy sport rules can be set and changed within a
// league; for instance, the roster positions, statistics used to score,
// scoring modifiers, and game style are configurable.
//
// The game structure means that a lot of fantasy data is relevant only in the
// context of a particular league and team. For instance, without the league’s
// scoring rules, the statistics compiled by a player in a real-life competition
// are not meaningful to a particular league. Three rushing touchdowns by a
// running back is irrelevant to a league that only considers defensive players.
// Many leagues are private – the information about them is only available to
// users that are a members.
//
//
//OAuth
//
// If you’re going to use the Fantasy Sports APIs, you’re going to have to get a
// bit familiar with OAuth. OAuth is the authentication mechanism for these
// services that allows users to grant you permission to make requests on their
// behalf. Many other Yahoo! services use OAuth, and thus all of the underlying
// details are explained in exhaustive detail in our primary OAuth
// documentation. Of particular interest is the OAuth Authorization Flow, which
// explains where each request is made and where the user needs to get involved.
//
// However, constructing OAuth flows from scratch is complicated and easy to
// get wrong. It’s often easier to use existing libraries, which are available
// for most languages on the OAuth.net Code page.
//
//
// Registering Your Application
//
// To work with OAuth and Yahoo! services, you also must register your
// application with the Yahoo! Developer Network. When you register your
// application, you define a scope of Yahoo! services that your application
// will need access to, as well as the basic descriptive information that will
// be presented to users of your application when they’re asked to grant you
// permissions. You will be given a consumer key and secret value that will
// need to be fed into OAuth requests that you generate. You should be sure to
// keep these values secret, as anyone with access to them could masquerade as
// your application.
//
// To create a new OAuth application to use with the Fantasy Sports APIs, you
// should go through the New API Key flow on YDN. Be sure to specify that you
// need access to private user data, and select either Read or Read/Write access
// for Fantasy Sports.
//
//
// Resources and Collections
//
// Introduction
// The primary building blocks of the Fantasy Sports APIs are Resources and
// Collections. Resources typically describe chunks of data that can be
// identified by a unique key. Collections are simply wrappers that contain
// similar resources. So, for instance, if we need to retrieve data about a
// single league, we might ask for a League Resource and provide a single league
// key. However, if we wanted data across several leagues, we would ask for a
// Leagues Collection and provide multiple league keys.
//
// The format for requesting a Resource will typically look like:
//
// http://fantasysports.yahooapis.com/fantasy/v2/{resource}/{resource_key}
//
// While the format for requesting a Collection will typically look like:
//
// http://fantasysports.yahooapis.com/fantasy/v2/{collection};{resource}_keys={resource_key1},{resource_key2}
//
// Collections
// As mentioned, Collections are simply groups of Resources. If you care about
// particular Resources within a Collection, you can apply filters to the
// Collection to narrow the results. The most common type of filtering is by
// key. For instance, if you’d like to see two particular players, you could ask
// for them directly:
//
// /fantasy/v2/;player_keys=,{player_key2}
//
// Some Collections support more complex filters. Within a game, for example,
// you might ask for only the players that play a certain position. You could
// also request only the particular user who is currently logged in:
//
// /fantasy/v2/;use_login=1
//
// Sub-Resources
// Resources will typically define a list of valid Sub-Resources. These are
// Resources and Collections that can live within the scope of the parent
// Resource. For instance, a fantasy league in the Football draft and trade
// games can contain up to 20 fantasy teams; therefore, the League Resource can
// have a Teams Collection as a sub-resource. As a general rule of thumb, if you
// can possibly have multiple of one Resource contained within another Resource,
// then you’ll have a Collection of the first Resource as a sub-resource of the
// second Resource. You would only have a singular Resource as a sub-resource of
// another Resource if it would only make sense to ever ask for one instance of
// that Resource.
//
// The scope of a sub-resource is typically defined by the parent Resource; for
// instance, if you’re viewing a Players Collection as a sub-resource of a
// particular League, then you would expect to only see Players that are
// eligible within that League. Further filters could then be applied to this
// already narrowed list.
//
// Having sub-resources allows you to chain together Resources and Collections
// to provide more data, and the URI you request directly specifies how the
// chaining works. For instance, if you wanted to take a particular logged in
// user, see which games he had played, and then get the league information
// within those games, you might construct a request like:
//   /fantasy/v2/;use_login=1//
//
// This would present you with a Users Collection, a single User Resource for
// the logged in user, a Games Collection for that user, potentially multiple
// Game Resources for each game the user is playing, a Leagues Collection
// beneath each Game Resource, and potentially multiple League Resources for
// each league the user belongs to in that game.
//
// When you specify a sub-resource beneath a Collection, you’re really saying
// that you’d like to see that sub-resource appended beneath each Resource
// within the Collection. Therefore, the sub-resources available to a
// Collection will be equivalent to the sub-resources available to the
// corresponding Resource.
//
// If you ever need to branch off other sub-resources outside of your main
// resource chain, you can use the out parameter, which will let you specify
// one level of extra sub-resources to pull in. At the moment, you cannot pass
// any parameters along to these out sub-resources, aside from any data that
// might get passed by default. This typically means that you can’t chain other
// resources off of sub-resources specified by the out parameter.
//
// As an example, if you wanted to view a league’s settings along with two teams
// in particular in a league, you might construct a URI like:
//
//   /fantasy/v2//;out=settings/;team_keys=,{team_key2}
//
// Parameters
// Parameters can be provided to Resources and Collections as semicolon-
// delimited key-value pairs. These should be placed after the Resource or
// Collection name in the URI; in the case of entry-point Resources like Games,
// Leagues, Teams, and Players, the parameters belong after the resource_key.
//   /fantasy/v2/{resource}/{resource_key};{key}={value};{key}={value}/{collection};{key}={value};{key}={value}/{collection};{key}={value}/{resource};{key}={value}
//
// Resource keys, out parameters, and other filters are just specific types of
// parameters that can be applied to various Resources or Collections.
//

// Game resource
//
// Description
// With the Game API, you can obtain the fantasy game related information, like
// the fantasy game name, the Yahoo! game code, and season.
//
// To refer to a Game resource, you’ll need to provide a game_key, which will
// either be a game_id or game_code. The game_id is a unique ID identifying a
// given fantasy game for a given season. For instance, the game_id for the Free
// NFL draft and trade fantasy game for the 2009 season is 222, while the
// game_id for the Plus version is 223. A game_code generally identifies a game,
// independent of season, and, when used as a game_key, will typically return
// the current season of that game. For instance, the game_code for the Free NFL
// game is nfl, and the game_code for the Plus game is pnfl; using nfl as your
// game_key during the 2010 season would be the same as providing the game_id
// for the 2010 season of the NFL game (242). As of the 2010 seasons, the Plus
// and Free games have been combined into a single code. Next year, the
// game_code nfl will point to the new game_id for the 2011 version of the NFL
// game. Thus, if you always want the current season of a game, the game_code
// should be used as a game_key.
//
// Below is a list of game IDs for most of our seasons of each game. If you’re
// looking for a current game ID that’s not listed in the table below, as
// mentioned above, you can request game information by game_code. For example:
//
// YQL: select * from fantasysports.games where game_key='nfl';
// API: http://fantasysports.yahooapis.com/fantasy/v2/game/nfl
// Game IDs Table
// Season  nfl game ID pnfl game ID  mlb game ID pmlb game ID  nba game ID  nhl game ID
// 2001    57          58            12          –             16           15
// 2002    49          62            39          44            67           64
// 2003    79          78            74          73            95           94
// 2004    101         102           98          99            112          111
// 2005    124         125           113         114           131          130
// 2006    153         154           147         148           165          164
// 2007    175         176           171         172           187          186
// 2008    199         200           195         196           211          210
// 2009    222         223           215         216           234          233
// 2010    242         –             238         –             249          248
// 2011    257         –             253         –             265          263
// 2012    273         –             268         –             304          303
//
//    <?xml version="1.0" encoding="UTF-8"?>
//     <fantasy_content xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/game/nfl" xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" time="30.575037002563ms" copyright="Data provided by Yahoo! and STATS, LLC" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
//       <game>
//         <game_key>257</game_key>
//         <game_id>257</game_id>
//         <name>Football</name>
//         <code>nfl</code>
//         <type>full</type>
//         <url>http://football.fantasysports.yahoo.com/f1</url>
//         <season>2011</season>
//       </game>
//     </fantasy_content>
//
type GameResource struct {
	XMLName  xml.Name `xml:"game"`
	Game_key string   `xml:"game_key"`
	Game_id  string   `xml:"game_id"`
	Namecode string   `xml:"namecode"`
	Type     string   `xml:"type"`
	Url      string   `xml:"url"`
	Season   string   `xml:"season"`
}

/*
HTTP Operations Supported

GET
URIs¶

http://fantasysports.yahooapis.com/fantasy/v2/game/

Any sub-resource under a game is extracted using a URI like:

http://fantasysports.yahooapis.com/fantasy/v2/game//

Multiple sub-resources can be extracted from game in the same URI using a format like:

http://fantasysports.yahooapis.com/fantasy/v2/game/;out=,{sub_resource_2}

Game key format¶

{game_code} or {game_id}

Example:pnfl or 223

Note

If you specify a game_code as the game_key , we’ll translate that to the corresponding game_id upon parsing the URI. Therefore, any game_code s will be converted to game_id s in any keys returned by the Fantasy Sports APIs in the response XML.

Sub-resources¶

Default sub-resource: metadata

Name	Description	URI	Sample
metadata	Includes game key, code, name, url, type and season.	/fantasy/v2/game//metadata	The 2009 Football PLUS game: http://fantasysports.yahooapis.com/fantasy/v2/game/223
 	Fetch specified leagues under a game.	/fantasy/v2/game//leagues;league_keys=,{league_key2}	A publicly viewable league within the 2009 football plus game: http://fantasysports.yahooapis.com/fantasy/v2/game/223/leagues;league_keys=223.l.431
 	Fetch specified players under a game.	/fantasy/v2/game//players;player_keys=,{player_key2}	Brett Favre’s information from the 2009 football plus game: http://fantasysports.yahooapis.com/fantasy/v2/game/223/players;player_keys=223.p.1025
game_weeks	Start and end date information for each week in the game	/fantasy/v2/game//game_weeks	NFL game weeks http://fantasysports.yahooapis.com/fantasy/v2/game/nfl/game_weeks
stat_categories	Detailed description of all available stat categories for the game.	/fantasy/v2/game//stat_categories	NFL stat categories http://fantasysports.yahooapis.com/fantasy/v2/game/nfl/stat_categories
position_types	Detailed description of all player position types for the game.	/fantasy/v2/game//position_types	NFL position types http://fantasysports.yahooapis.com/fantasy/v2/game/nfl/position_types
roster_positions	Detailed description of all roster positions for the game.	/fantasy/v2/game//roster_positions	NFL roster positions http://fantasysports.yahooapis.com/fantasy/v2/game/nfl/roster_positions
Sample XML¶

http://fantasysports.yahooapis.com/fantasy/v2/game/nfl

*/

/*
Games collection¶

Description¶

With the Games API, you can obtain information from a collection of games simultaneously. Each element beneath the Games Collection will be a Game Resource

HTTP Operations Supported¶

GET
URIs¶

URI	Description	Sample
http://fantasysports.yahooapis.com/fantasy/v2/games;game_keys=,{game_key2}	Fetch specific games {game_key1} and {game_key2}	nfl and mlb games: http://fantasysports.yahooapis.com/fantasy/v2/games;game_keys=nfl,mlb
http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games	Fetch all games for the logged in user	all games for user: http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games
http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys={game_key1},{game_key2}	Fetch specific games {game_key1} and {game_key2} that the logged in user owns teams in.	nfl and mlb games for user: http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys=nfl,mlb
Any sub-resource valid for a game is a valid sub-resource under the games collection.

Any sub-resource for a collection of games is extracted using a URI like:

/games/{sub_resource}

OR

/games;game_keys={game_key1},{game_key2}/{sub_resource}

Multiple sub-resources can be extracted from games in the same URI using a format like:

/games;out={sub_resource_1},{sub_resource_2}

OR

/games;game_keys={game_key1},{game_key2};out={sub_resource_1},{sub_resource_2}

Filters¶

The games collection can have filters such as the following to obtain a subset of a games collection that satisfy the filtering condition. These filters can be combined to obtain a more restricted list of games. For instance, if you wanted only the 2011 version of the nfl game, you might filter by seasons=2011 and game_codes=nfl.

Filter parameter	Filter parameter values	Usage
is_available	1 to only show games currently in season	/games;is_available=1
game_types	full|pickem-team|pickem-group|pickem-team-list	/games;game_types=full,pickem-team
game_codes	Any valid game codes	/games;game_codes=nfl,mlb
seasons	Any valid seasons	/games;seasons=2011,2012
Sub-resources¶

In addition to the sub-resources valid for a game resource, the following are valid sub-resources for a games collection.

Name	Description	URI	Sample
 	Fetch teams owned by a user for one or more games.
/fantasy/v2/users;use_login=1/games/teams

OR

/fantasy/v2/users;use_login=1/games;game_keys={game_key1},{game_key2}/teams

all teams for user: http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games/teams
*/

/*
League resource¶

Description¶

When users join a Fantasy Football, Baseball, Basketball, or Hockey draft and trade game, they are organized into leagues with a limited number of friends or other Yahoo! users, with each user managing a Team. With the League API, you can obtain the league related information, like the league name, the number of teams, the draft status, et cetera. Leagues only exist in the context of a particular Game, although you can request a League Resource as the base of your URI by using the global ````. A particular user can only retrieve data for private leagues of which they are a member, or for public leagues.

HTTP Operations Supported¶

GET
URIs¶

http://fantasysports.yahooapis.com/fantasy/v2/league/

Any sub-resource under a league is extracted using a URI like:

http://fantasysports.yahooapis.com/fantasy/v2/league//

Multiple sub-resources can be extracted from league in the same URI using a format like:

http://fantasysports.yahooapis.com/fantasy/v2/league/;out=,{sub_resource_2}

League key format¶

.l.{league_id}

Example:pnfl.l.431 or 223.l.431

Note

The separator between the game_key and league_id is a lower case L (not the number 1).

Sub-resources¶

Default sub-resource: metadata

Name	Description	URI	Sample
metadata	Includes league key, id, name, url, draft status, number of teams, and current week information.	/fantasy/v2/league//metadata	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431
settings	League settings. For instance, draft type, scoring type, roster positions, stat categories and modifiers, divisions.	/fantasy/v2/league//settings	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/settings
standings	Ranking of teams within the league. Accepts Teams as a sub-resource, and includes team_standings data by default beneath the teams	/fantasy/v2/league//standings	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/standings
scoreboard	League scoreboard. Accepts Matchups as a sub-resource, which in turn accept Teams as a sub-resource. Includes team_stats data by default.
Scoreboard for current week: /fantasy/v2/league//scoreboard

Scoreboard for a particular week: /fantasy/v2/league//scoreboard;week={week}

http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/scoreboard;week=2
``

``
All teams in the league.	/fantasy/v2/league//teams	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/teams
``

``
The league’s eligible players.	/fantasy/v2/league//players	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players
draftresults	Draft results for all teams in the league.	/fantasy/v2/league//draftresults	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/draftresults
``

``
League transactions – adds, drops, and trades.	/fantasy/v2/league//transactions	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/transactions
Sample XML¶

http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431" xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" time="181.80584907532ms" copyright="Data provided by Yahoo! and STATS, LLC" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
  <league>
    <league_key>223.l.431</league_key>
    <league_id>431</league_id>
    <name>Y! Friends and Family League</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431</url>
    <draft_status>postdraft</draft_status>
    <num_teams>14</num_teams>
    <edit_key>17</edit_key>
    <weekly_deadline/>
    <league_update_timestamp>1262595518</league_update_timestamp>
    <scoring_type>head</scoring_type>
    <current_week>16</current_week>
    <start_week>1</start_week>
    <end_week>16</end_week>
    <is_finished>1</is_finished>
  </league>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/settings

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/settings" time="86.472988128662ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <league>
    <league_key>223.l.431</league_key>
    <league_id>431</league_id>
    <name>Y! Friends and Family League</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431</url>
    <draft_status>postdraft</draft_status>
    <num_teams>14</num_teams>
    <edit_key>17</edit_key>
    <weekly_deadline/>
    <league_update_timestamp>1262595518</league_update_timestamp>
    <scoring_type>head</scoring_type>
    <current_week>16</current_week>
    <start_week>1</start_week>
    <end_week>16</end_week>
    <is_finished>1</is_finished>
    <settings>
      <draft_type>live</draft_type>
      <scoring_type>head</scoring_type>
      <uses_playoff>1</uses_playoff>
      <playoff_start_week>14</playoff_start_week>
      <uses_playoff_reseeding>0</uses_playoff_reseeding>
      <uses_lock_eliminated_teams>0</uses_lock_eliminated_teams>
      <uses_faab>1</uses_faab>
      <trade_end_date>2009-11-27</trade_end_date>
      <trade_ratify_type>commish</trade_ratify_type>
      <trade_reject_time>0</trade_reject_time>
      <roster_positions>
        <roster_position>
          <position>QB</position>
          <count>1</count>
        </roster_position>
        <roster_position>
          <position>WR</position>
          <count>3</count>
        </roster_position>
        <roster_position>
          <position>RB</position>
          <count>2</count>
        </roster_position>
        <roster_position>
          <position>TE</position>
          <count>1</count>
        </roster_position>
        <roster_position>
          <position>W/R/T</position>
          <count>1</count>
        </roster_position>
        <roster_position>
          <position>K</position>
          <count>1</count>
        </roster_position>
        <roster_position>
          <position>DEF</position>
          <count>1</count>
        </roster_position>
        <roster_position>
          <position>BN</position>
          <count>4</count>
        </roster_position>
      </roster_positions>
      <stat_categories>
        <stats>
          <stat>
            <stat_id>4</stat_id>
            <enabled>1</enabled>
            <name>Passing Yards</name>
            <display_name>Pass Yds</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>5</stat_id>
            <enabled>1</enabled>
            <name>Passing Touchdowns</name>
            <display_name>Pass TD</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>6</stat_id>
            <enabled>1</enabled>
            <name>Interceptions</name>
            <display_name>Int</display_name>
            <sort_order>0</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>9</stat_id>
            <enabled>1</enabled>
            <name>Rushing Yards</name>
            <display_name>Rush Yds</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>10</stat_id>
            <enabled>1</enabled>
            <name>Rushing Touchdowns</name>
            <display_name>Rush TD</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>11</stat_id>
            <enabled>1</enabled>
            <name>Receptions</name>
            <display_name>Rec</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>12</stat_id>
            <enabled>1</enabled>
            <name>Reception Yards</name>
            <display_name>Rec Yds</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>13</stat_id>
            <enabled>1</enabled>
            <name>Reception Touchdowns</name>
            <display_name>Rec TD</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>15</stat_id>
            <enabled>1</enabled>
            <name>Return Touchdowns</name>
            <display_name>Ret TD</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>16</stat_id>
            <enabled>1</enabled>
            <name>2-Point Conversions</name>
            <display_name>2-PT</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>18</stat_id>
            <enabled>1</enabled>
            <name>Fumbles Lost</name>
            <display_name>Fum Lost</display_name>
            <sort_order>0</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>57</stat_id>
            <enabled>1</enabled>
            <name>Offensive Fumble Return TD</name>
            <display_name>Fum Ret TD</display_name>
            <sort_order>1</sort_order>
            <position_type>O</position_type>
          </stat>
          <stat>
            <stat_id>19</stat_id>
            <enabled>1</enabled>
            <name>Field Goals 0-19 Yards</name>
            <display_name>FG 0-19</display_name>
            <sort_order>1</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>20</stat_id>
            <enabled>1</enabled>
            <name>Field Goals 20-29 Yards</name>
            <display_name>FG 20-29</display_name>
            <sort_order>1</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>21</stat_id>
            <enabled>1</enabled>
            <name>Field Goals 30-39 Yards</name>
            <display_name>FG 30-39</display_name>
            <sort_order>1</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>22</stat_id>
            <enabled>1</enabled>
            <name>Field Goals 40-49 Yards</name>
            <display_name>FG 40-49</display_name>
            <sort_order>1</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>23</stat_id>
            <enabled>1</enabled>
            <name>Field Goals 50+ Yards</name>
            <display_name>FG 50+</display_name>
            <sort_order>1</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>24</stat_id>
            <enabled>1</enabled>
            <name>Field Goals Missed 0-19 Yards</name>
            <display_name>FGM 0-19</display_name>
            <sort_order>0</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>25</stat_id>
            <enabled>1</enabled>
            <name>Field Goals Missed 20-29 Yards</name>
            <display_name>FGM 20-29</display_name>
            <sort_order>0</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>29</stat_id>
            <enabled>1</enabled>
            <name>Point After Attempt Made</name>
            <display_name>PAT Made</display_name>
            <sort_order>1</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>30</stat_id>
            <enabled>1</enabled>
            <name>Point After Attempt Missed</name>
            <display_name>PAT Miss</display_name>
            <sort_order>0</sort_order>
            <position_type>K</position_type>
          </stat>
          <stat>
            <stat_id>31</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed</name>
            <display_name>Pts Allow</display_name>
            <sort_order>0</sort_order>
            <position_type>DT</position_type>
            <is_only_display_stat>1</is_only_display_stat>
          </stat>
          <stat>
            <stat_id>32</stat_id>
            <enabled>1</enabled>
            <name>Sack</name>
            <display_name>Sack</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>33</stat_id>
            <enabled>1</enabled>
            <name>Interception</name>
            <display_name>Int</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>34</stat_id>
            <enabled>1</enabled>
            <name>Fumble Recovery</name>
            <display_name>Fum Rec</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>35</stat_id>
            <enabled>1</enabled>
            <name>Touchdown</name>
            <display_name>TD</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>36</stat_id>
            <enabled>1</enabled>
            <name>Safety</name>
            <display_name>Safe</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>37</stat_id>
            <enabled>1</enabled>
            <name>Block Kick</name>
            <display_name>Blk Kick</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>50</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 0 points</name>
            <display_name>Pts Allow 0</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>51</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 1-6 points</name>
            <display_name>Pts Allow 1-6</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>52</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 7-13 points</name>
            <display_name>Pts Allow 7-13</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>53</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 14-20 points</name>
            <display_name>Pts Allow 14-20</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>54</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 21-27 points</name>
            <display_name>Pts Allow 21-27</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>55</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 28-34 points</name>
            <display_name>Pts Allow 28-34</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
          <stat>
            <stat_id>56</stat_id>
            <enabled>1</enabled>
            <name>Points Allowed 35+ points</name>
            <display_name>Pts Allow 35+</display_name>
            <sort_order>1</sort_order>
            <position_type>DT</position_type>
          </stat>
        </stats>
      </stat_categories>
      <stat_modifiers>
        <stats>
          <stat>
            <stat_id>4</stat_id>
            <value>0.04</value>
          </stat>
          <stat>
            <stat_id>5</stat_id>
            <value>4</value>
          </stat>
          <stat>
            <stat_id>6</stat_id>
            <value>-1</value>
          </stat>
          <stat>
            <stat_id>9</stat_id>
            <value>0.1</value>
          </stat>
          <stat>
            <stat_id>10</stat_id>
            <value>6</value>
          </stat>
          <stat>
            <stat_id>11</stat_id>
            <value>.75</value>
          </stat>
          <stat>
            <stat_id>12</stat_id>
            <value>0.1</value>
          </stat>
          <stat>
            <stat_id>13</stat_id>
            <value>6</value>
          </stat>
          <stat>
            <stat_id>15</stat_id>
            <value>6</value>
          </stat>
          <stat>
            <stat_id>16</stat_id>
            <value>2</value>
          </stat>
          <stat>
            <stat_id>18</stat_id>
            <value>-1</value>
          </stat>
          <stat>
            <stat_id>57</stat_id>
            <value>6</value>
          </stat>
          <stat>
            <stat_id>19</stat_id>
            <value>3</value>
          </stat>
          <stat>
            <stat_id>20</stat_id>
            <value>3</value>
          </stat>
          <stat>
            <stat_id>21</stat_id>
            <value>3</value>
          </stat>
          <stat>
            <stat_id>22</stat_id>
            <value>4</value>
          </stat>
          <stat>
            <stat_id>23</stat_id>
            <value>5</value>
          </stat>
          <stat>
            <stat_id>24</stat_id>
            <value>-3</value>
          </stat>
          <stat>
            <stat_id>25</stat_id>
            <value>-1</value>
          </stat>
          <stat>
            <stat_id>29</stat_id>
            <value>1</value>
          </stat>
          <stat>
            <stat_id>30</stat_id>
            <value>-.5</value>
          </stat>
          <stat>
            <stat_id>32</stat_id>
            <value>1</value>
          </stat>
          <stat>
            <stat_id>33</stat_id>
            <value>2</value>
          </stat>
          <stat>
            <stat_id>34</stat_id>
            <value>2</value>
          </stat>
          <stat>
            <stat_id>35</stat_id>
            <value>6</value>
          </stat>
          <stat>
            <stat_id>36</stat_id>
            <value>2</value>
          </stat>
          <stat>
            <stat_id>37</stat_id>
            <value>2</value>
          </stat>
          <stat>
            <stat_id>50</stat_id>
            <value>10</value>
          </stat>
          <stat>
            <stat_id>51</stat_id>
            <value>7</value>
          </stat>
          <stat>
            <stat_id>52</stat_id>
            <value>4</value>
          </stat>
          <stat>
            <stat_id>53</stat_id>
            <value>1</value>
          </stat>
          <stat>
            <stat_id>54</stat_id>
            <value>0</value>
          </stat>
          <stat>
            <stat_id>55</stat_id>
            <value>-1</value>
          </stat>
          <stat>
            <stat_id>56</stat_id>
            <value>-4</value>
          </stat>
        </stats>
      </stat_modifiers>
      <divisions>
        <division>
          <division_id>1</division_id>
          <name>Family</name>
        </division>
        <division>
          <division_id>2</division_id>
          <name>Friends</name>
        </division>
      </divisions>
    </settings>
  </league>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/standings

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/standings" time="201.46489143372ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <league>
    <league_key>223.l.431</league_key>
    <league_id>431</league_id>
    <name>Y! Friends and Family League</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431</url>
    <draft_status>postdraft</draft_status>
    <num_teams>14</num_teams>
    <edit_key>17</edit_key>
    <weekly_deadline/>
    <league_update_timestamp>1262595518</league_update_timestamp>
    <scoring_type>head</scoring_type>
    <current_week>16</current_week>
    <start_week>1</start_week>
    <end_week>16</end_week>
    <is_finished>1</is_finished>
    <standings>
      <teams count="14">
        <team>
          <team_key>223.l.431.t.10</team_key>
          <team_id>10</team_id>
          <name>Gehlken</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/10</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://a323.yahoofs.com/coreid/4b978f0ci2432zws140sp2/imXqmYo8cq3NxEFtQB4wgAs-/6/tn48.jpeg?ciA8DVOBMH.UXGXk</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>0</faab_balance>
          <clinched_playoffs>1</clinched_playoffs>
          <managers>
            <manager>
              <manager_id>5</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>4LAITFUXFASDNAXFWUOHWNU3BY</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1682.33</total>
          </team_points>
          <team_standings>
            <rank>1</rank>
            <outcome_totals>
              <wins>9</wins>
              <losses>4</losses>
              <ties>0</ties>
              <percentage>.692</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>5</wins>
              <losses>1</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.5</team_key>
          <team_id>5</team_id>
          <name>RotoExperts</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/5</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://a323.yahoofs.com/coreid/49be42a6i26e5zul3re3/d2x_9_UweKP95SJZ_Hwnk2Rl/2/tn48.jpg?ciA8DVOBIRa6b7wq</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>1</faab_balance>
          <clinched_playoffs>1</clinched_playoffs>
          <managers>
            <manager>
              <manager_id>12</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>RW3ELDFMOFTES2EUAWQVCPPN7E</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1764.09</total>
          </team_points>
          <team_standings>
            <rank>2</rank>
            <outcome_totals>
              <wins>9</wins>
              <losses>4</losses>
              <ties>0</ties>
              <percentage>.692</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>4</wins>
              <losses>2</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.8</team_key>
          <team_id>8</team_id>
          <name>Y! - Pianowski</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/8</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_10_48.gif</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>0</faab_balance>
          <clinched_playoffs>1</clinched_playoffs>
          <managers>
            <manager>
              <manager_id>6</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>WMKEJTV3VUJA4VZWQ25O27W43M</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1569.48</total>
          </team_points>
          <team_standings>
            <rank>3</rank>
            <outcome_totals>
              <wins>8</wins>
              <losses>5</losses>
              <ties>0</ties>
              <percentage>.615</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>4</wins>
              <losses>2</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.12</team_key>
          <team_id>12</team_id>
          <name>Y! - Behrens</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/12</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://lookup.avatars.yahoo.com/images?yid=abehrens53&amp;size=medium&amp;type=jpg&amp;pty=3000</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>0</faab_balance>
          <clinched_playoffs>1</clinched_playoffs>
          <managers>
            <manager>
              <manager_id>3</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>E2KS77CDQPACRTSBCYPOFFW6AI</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1652.27</total>
          </team_points>
          <team_standings>
            <rank>4</rank>
            <outcome_totals>
              <wins>8</wins>
              <losses>5</losses>
              <ties>0</ties>
              <percentage>.615</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>5</wins>
              <losses>1</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.4</team_key>
          <team_id>4</team_id>
          <name>Salfino-Comcast/NESN</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/4</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://a323.yahoofs.com/coreid/4d8a517fi1b71zul1re3/ypdMGIA8cbVafvybuj2J.Jg-/2/tn48.jpg?ciA8DVOB5bipYD0R</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>0</faab_balance>
          <clinched_playoffs>1</clinched_playoffs>
          <managers>
            <manager>
              <manager_id>9</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>PDLVXDDVXK2FRDI3FHRSS74F2U</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1621.98</total>
          </team_points>
          <team_standings>
            <rank>5</rank>
            <outcome_totals>
              <wins>7</wins>
              <losses>6</losses>
              <ties>0</ties>
              <percentage>.538</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>3</wins>
              <losses>3</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.11</team_key>
          <team_id>11</team_id>
          <name>FantasyGuru.com-Hans</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/11</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://lookup.avatars.yahoo.com/images?yid=fantasygurudotcom&amp;size=medium&amp;type=jpg&amp;pty=3000</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>1</faab_balance>
          <clinched_playoffs>1</clinched_playoffs>
          <managers>
            <manager>
              <manager_id>4</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>B7IJFDI5UUTN3AQ2F7ZEA4BDU4</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1469.00</total>
          </team_points>
          <team_standings>
            <rank>6</rank>
            <outcome_totals>
              <wins>7</wins>
              <losses>6</losses>
              <ties>0</ties>
              <percentage>.538</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>2</wins>
              <losses>4</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.1</team_key>
          <team_id>1</team_id>
          <name>PFW - Blunda</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/1</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_01_48.gif</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>22</faab_balance>
          <managers>
            <manager>
              <manager_id>13</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>XNAXQZRDZPJ3RVFMY7ZTSWEFLU</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1461.71</total>
          </team_points>
          <team_standings>
            <rank>7</rank>
            <outcome_totals>
              <wins>7</wins>
              <losses>6</losses>
              <ties>0</ties>
              <percentage>.538</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>3</wins>
              <losses>3</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.2</team_key>
          <team_id>2</team_id>
          <name>Y! - Evans</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/2</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://a323.yahoofs.com/coreid/4a68b2d6i2663zul3re3/HYebAP0zcqEPfMp3gOK8Mmbv/4/tn48.jpg?ciA8DVOBzMjxdtsK</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>35</faab_balance>
          <managers>
            <manager>
              <manager_id>8</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>RV2NLFT5LDNKUDOFSWSHIDINY4</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1512.53</total>
          </team_points>
          <team_standings>
            <rank>8</rank>
            <outcome_totals>
              <wins>6</wins>
              <losses>7</losses>
              <ties>0</ties>
              <percentage>.462</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>2</wins>
              <losses>4</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.13</team_key>
          <team_id>13</team_id>
          <name>Erickson - RotoWire</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/13</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://lookup.avatars.yahoo.com/images?yid=jeff_rotonews&amp;size=medium&amp;type=jpg&amp;pty=3000</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>17</faab_balance>
          <managers>
            <manager>
              <manager_id>11</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>SB4Y5HVVUKMCTKZFQCXHIZ222E</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1484.56</total>
          </team_points>
          <team_standings>
            <rank>9</rank>
            <outcome_totals>
              <wins>6</wins>
              <losses>7</losses>
              <ties>0</ties>
              <percentage>.462</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>3</wins>
              <losses>3</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.9</team_key>
          <team_id>9</team_id>
          <name>Y! - Funston</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/9</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://lookup.avatars.yahoo.com/images?yid=brandoanf1&amp;size=medium&amp;type=jpg&amp;pty=3000</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>10</faab_balance>
          <managers>
            <manager>
              <manager_id>1</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>3H7IQ3F2742K2ODHSJK5YXL23E</guid>
              <is_commissioner>1</is_commissioner>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1430.24</total>
          </team_points>
          <team_standings>
            <rank>10</rank>
            <outcome_totals>
              <wins>6</wins>
              <losses>7</losses>
              <ties>0</ties>
              <percentage>.462</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>2</wins>
              <losses>4</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.7</team_key>
          <team_id>7</team_id>
          <name>RotoWire_Liss</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/7</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_10_48.gif</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>68</faab_balance>
          <managers>
            <manager>
              <manager_id>7</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>4BDB5LIG3IFVROH7SRBX44LBZM</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1424.56</total>
          </team_points>
          <team_standings>
            <rank>11</rank>
            <outcome_totals>
              <wins>6</wins>
              <losses>7</losses>
              <ties>0</ties>
              <percentage>.462</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>3</wins>
              <losses>3</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.3</team_key>
          <team_id>3</team_id>
          <name>RotoWire - Del Don</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/3</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_05_48.gif</url>
            </team_logo>
          </team_logos>
          <division_id>2</division_id>
          <faab_balance>0</faab_balance>
          <managers>
            <manager>
              <manager_id>10</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>4A5KVYHC7ZSEGOBFHFSO5Q64VA</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1366.89</total>
          </team_points>
          <team_standings>
            <rank>12</rank>
            <outcome_totals>
              <wins>6</wins>
              <losses>7</losses>
              <ties>0</ties>
              <percentage>.462</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>3</wins>
              <losses>3</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.6</team_key>
          <team_id>6</team_id>
          <name>Y! - Romig</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/6</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://a323.yahoofs.com/coreid/49b954dci229az/IJtbcRQjdKtd_DMoStSK/103/tn48.jpg?ciA8DVOB2WQ2Fk4F</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>0</faab_balance>
          <managers>
            <manager>
              <manager_id>2</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>FS5M5LOFJRKVJNRIWG36ZUF7IQ</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1370.16</total>
          </team_points>
          <team_standings>
            <rank>13</rank>
            <outcome_totals>
              <wins>5</wins>
              <losses>8</losses>
              <ties>0</ties>
              <percentage>.385</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>2</wins>
              <losses>4</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
        <team>
          <team_key>223.l.431.t.14</team_key>
          <team_id>14</team_id>
          <name>Y! - Chase</name>
          <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/14</url>
          <team_logos>
            <team_logo>
              <size>medium</size>
              <url>http://a323.yahoofs.com/coreid/4a7a23a5icfazul2re3/2fIcrk8yc7QS3j_ei4PULEbpFA--/1/tn48.jpg?ciA8DVOBcEQk3vWZ</url>
            </team_logo>
          </team_logos>
          <division_id>1</division_id>
          <faab_balance>92</faab_balance>
          <managers>
            <manager>
              <manager_id>14</manager_id>
              <nickname>-- hidden --</nickname>
              <guid>7CSOKBMM74MGFMSWHWJMM4FBQ4</guid>
            </manager>
          </managers>
          <team_points>
            <coverage_type>season</coverage_type>
            <season>2009</season>
            <total>1237.47</total>
          </team_points>
          <team_standings>
            <rank>14</rank>
            <outcome_totals>
              <wins>1</wins>
              <losses>12</losses>
              <ties>0</ties>
              <percentage>.077</percentage>
            </outcome_totals>
            <divisional_outcome_totals>
              <wins>1</wins>
              <losses>5</losses>
              <ties>0</ties>
            </divisional_outcome_totals>
          </team_standings>
        </team>
      </teams>
    </standings>
  </league>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/scoreboard

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/scoreboard" time="148.71311187744ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <league>
    <league_key>223.l.431</league_key>
    <league_id>431</league_id>
    <name>Y! Friends and Family League</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431</url>
    <draft_status>postdraft</draft_status>
    <num_teams>14</num_teams>
    <edit_key>17</edit_key>
    <weekly_deadline/>
    <league_update_timestamp>1262595518</league_update_timestamp>
    <scoring_type>head</scoring_type>
    <current_week>16</current_week>
    <start_week>1</start_week>
    <end_week>16</end_week>
    <is_finished>1</is_finished>
    <scoreboard>
      <week>16</week>
      <matchups count="2">
        <matchup>
          <week>16</week>
          <status>postevent</status>
          <is_tied>0</is_tied>
          <winner_team_key>223.l.431.t.10</winner_team_key>
          <teams count="2">
            <team>
              <team_key>223.l.431.t.5</team_key>
              <team_id>5</team_id>
              <name>RotoExperts</name>
              <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/5</url>
              <team_logos>
                <team_logo>
                  <size>medium</size>
                  <url>http://a323.yahoofs.com/coreid/49be42a6i26e5zul3re3/d2x_9_UweKP95SJZ_Hwnk2Rl/2/tn48.jpg?ciA8DVOBIRa6b7wq</url>
                </team_logo>
              </team_logos>
              <division_id>2</division_id>
              <faab_balance>1</faab_balance>
              <clinched_playoffs>1</clinched_playoffs>
              <managers>
                <manager>
                  <manager_id>12</manager_id>
                  <nickname>-- hidden --</nickname>
                  <guid>RW3ELDFMOFTES2EUAWQVCPPN7E</guid>
                </manager>
              </managers>
              <team_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>135.22</total>
              </team_points>
              <team_projected_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>142.81</total>
              </team_projected_points>
            </team>
            <team>
              <team_key>223.l.431.t.10</team_key>
              <team_id>10</team_id>
              <name>Gehlken</name>
              <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/10</url>
              <team_logos>
                <team_logo>
                  <size>medium</size>
                  <url>http://a323.yahoofs.com/coreid/4b978f0ci2432zws140sp2/imXqmYo8cq3NxEFtQB4wgAs-/6/tn48.jpeg?ciA8DVOBMH.UXGXk</url>
                </team_logo>
              </team_logos>
              <division_id>1</division_id>
              <faab_balance>0</faab_balance>
              <clinched_playoffs>1</clinched_playoffs>
              <managers>
                <manager>
                  <manager_id>5</manager_id>
                  <nickname>-- hidden --</nickname>
                  <guid>4LAITFUXFASDNAXFWUOHWNU3BY</guid>
                </manager>
              </managers>
              <team_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>137.86</total>
              </team_points>
              <team_projected_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>133.57</total>
              </team_projected_points>
            </team>
          </teams>
        </matchup>
        <matchup>
          <week>16</week>
          <status>postevent</status>
          <is_tied>0</is_tied>
          <winner_team_key>223.l.431.t.8</winner_team_key>
          <teams count="2">
            <team>
              <team_key>223.l.431.t.8</team_key>
              <team_id>8</team_id>
              <name>Y! - Pianowski</name>
              <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/8</url>
              <team_logos>
                <team_logo>
                  <size>medium</size>
                  <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_10_48.gif</url>
                </team_logo>
              </team_logos>
              <division_id>1</division_id>
              <faab_balance>0</faab_balance>
              <clinched_playoffs>1</clinched_playoffs>
              <managers>
                <manager>
                  <manager_id>6</manager_id>
                  <nickname>-- hidden --</nickname>
                  <guid>WMKEJTV3VUJA4VZWQ25O27W43M</guid>
                </manager>
              </managers>
              <team_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>103.39</total>
              </team_points>
              <team_projected_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>104.17</total>
              </team_projected_points>
            </team>
            <team>
              <team_key>223.l.431.t.12</team_key>
              <team_id>12</team_id>
              <name>Y! - Behrens</name>
              <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/12</url>
              <team_logos>
                <team_logo>
                  <size>medium</size>
                  <url>http://lookup.avatars.yahoo.com/images?yid=abehrens53&amp;size=medium&amp;type=jpg&amp;pty=3000</url>
                </team_logo>
              </team_logos>
              <division_id>1</division_id>
              <faab_balance>0</faab_balance>
              <clinched_playoffs>1</clinched_playoffs>
              <managers>
                <manager>
                  <manager_id>3</manager_id>
                  <nickname>-- hidden --</nickname>
                  <guid>E2KS77CDQPACRTSBCYPOFFW6AI</guid>
                </manager>
              </managers>
              <team_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>101.94</total>
              </team_points>
              <team_projected_points>
                <coverage_type>week</coverage_type>
                <week>16</week>
                <total>127.28</total>
              </team_projected_points>
            </team>
          </teams>
        </matchup>
      </matchups>
    </scoreboard>
  </league>
</fantasy_content>
*/

/*
Leagues collection¶

Description¶

With the Leagues API, you can obtain information from a collection of leagues simultaneously. Each element beneath the Leagues Collection will be a League Resource

HTTP Operations Supported¶

GET
URIs¶

URI	Description	Sample
http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=,{league_key2}	Fetch specific leagues {league_key1} and {league_key2}	http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=223.l.431
Any sub-resource valid for a league is a valid sub-resource under the leagues collection.

Any sub-resource for a collection of leagues is extracted using a URI like:

/leagues/{sub_resource}

OR

/leagues;league_keys={league_key1},{league_key2}/{sub_resource}

Multiple sub-resources can be extracted from leagues in the same URI using a format like:

/leagues;out={sub_resource_1},{sub_resource_2}

OR

/leagues;league_keys={league_key1},{league_key2};out={sub_resource_1},{sub_resource_2}
*/

/*
Team resource¶

Description¶

The Team APIs allow you to retrieve information about a team within our fantasy games. The team is the basic unit for keeping track of a roster of players, and can be managed by either one or two managers (the second manager being called a co-manager). With the Team APIs, you can obtain team-related information, like the team name, managers, logos, stats and points, and rosters for particular weeks. Teams only exist in the context of a particular League, although you can request a Team Resource as the base of your URI by using the global ````. A particular user can only retrieve data about a team if that team is part of a private league of which the user is a member, or if it’s in a public league.

HTTP Operations Supported¶

GET
URIs¶

http://fantasysports.yahooapis.com/fantasy/v2/team/

Any sub-resource under a team is extracted using a URI like:

http://fantasysports.yahooapis.com/fantasy/v2/team//

Multiple sub-resources can be extracted from team in the same URI using a format like:

http://fantasysports.yahooapis.com/fantasy/v2/team/;out=,{sub_resource_2}

Team key format¶

.l.{league_id}.t.{team_id}

Example:pnfl.l.431.t.1 or 223.l.431.t.1

Sub-resources¶

Default sub-resource: metadata

Name	Description	URI	Sample
metadata	Includes team key, id, name, url, division ID, logos, and team manager information.	/fantasy/v2/team//metadata	http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9
stats	Team statistical data and points.
Season stats:/fantasy/v2/team//stats

Week stats: /fantasy/v2/team//stats;type=week;week={week}

Here {week} is a non-zero integer, or current for the current week.

http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9/stats;type=week;week=2
standings	Team rank, wins, losses, ties, and winning percentage (as well as divisional data if applicable).	/fantasy/v2/team//standings	http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9/standings
 	Team roster. Accepts a week parameter. Also accepts Players as a sub-resource (included by default)
Roster for a particular week: /fantasy/v2/team//roster;week={week}

Here {week} is a non-zero integer. If week is current, or isn’t provided, defaults to current week.

http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9/roster;week=2 - The week 2 roster for team 223.l.431.t.9
draftresults	List of players drafted by the team.	/fantasy/v2/team//draftresults	http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9/draftresults
matchups	All the matchups this team has scheduled (for H2H leagues).
All matchups: /fantasy/v2/team//matchups

Particular weeks: /fantasy/v2/team//matchups;weeks=1,3,6

http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9/matchups;weeks=1,3,6
Sample XML¶

http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1" time="426.26690864563ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <team>
    <team_key>223.l.431.t.1</team_key>
    <team_id>1</team_id>
    <name>PFW - Blunda</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/1</url>
    <team_logos>
      <team_logo>
        <size>medium</size>
        <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_01_48.gif</url>
      </team_logo>
    </team_logos>
    <division_id>2</division_id>
    <faab_balance>22</faab_balance>
    <managers>
      <manager>
        <manager_id>13</manager_id>
        <nickname>Michael Blunda</nickname>
        <guid>XNAXQZRDZPJ3RVFMY7ZTSWEFLU</guid>
      </manager>
    </managers>
  </team>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1/matchups;weeks=1,5 - team’s matchups for weeks 1 and 5 in a NFL H2H league

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1/matchups;weeks=1,5" time="576.54285430908ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <team>
    <team_key>223.l.431.t.1</team_key>
    <team_id>1</team_id>
    <name>PFW - Blunda</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/1</url>
    <team_logos>
      <team_logo>
        <size>medium</size>
        <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_01_48.gif</url>
      </team_logo>
    </team_logos>
    <division_id>2</division_id>
    <faab_balance>22</faab_balance>
    <managers>
      <manager>
        <manager_id>13</manager_id>
        <nickname>Michael Blunda</nickname>
        <guid>XNAXQZRDZPJ3RVFMY7ZTSWEFLU</guid>
      </manager>
    </managers>
    <matchups count="2">
      <matchup>
        <week>1</week>
        <status>postevent</status>
        <is_tied>0</is_tied>
        <winner_team_key>223.l.431.t.1</winner_team_key>
        <teams count="2">
          <team>
            <team_key>223.l.431.t.1</team_key>
            <team_id>1</team_id>
            <name>PFW - Blunda</name>
            <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/1</url>
            <team_logos>
              <team_logo>
                <size>medium</size>
                <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_01_48.gif</url>
              </team_logo>
            </team_logos>
            <division_id>2</division_id>
            <faab_balance>22</faab_balance>
            <managers>
              <manager>
                <manager_id>13</manager_id>
                <nickname>Michael Blunda</nickname>
                <guid>XNAXQZRDZPJ3RVFMY7ZTSWEFLU</guid>
              </manager>
            </managers>
            <team_points>
              <coverage_type>week</coverage_type>
              <week>1</week>
              <total>117.88</total>
            </team_points>
            <team_projected_points>
              <coverage_type>week</coverage_type>
              <week>1</week>
              <total>107.94</total>
            </team_projected_points>
          </team>
          <team>
            <team_key>223.l.431.t.5</team_key>
            <team_id>5</team_id>
            <name>RotoExperts</name>
            <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/5</url>
            <team_logos>
              <team_logo>
                <size>medium</size>
                <url>http://a323.yahoofs.com/coreid/49be42a6i26e5zul3re3/d2x_9_UweKP95SJZ_Hwnk2Rl/2/tn48.jpg?ciA8DVOBIRa6b7wq</url>
              </team_logo>
            </team_logos>
            <division_id>2</division_id>
            <faab_balance>1</faab_balance>
            <clinched_playoffs>1</clinched_playoffs>
            <managers>
              <manager>
                <manager_id>12</manager_id>
                <nickname>Scott</nickname>
                <guid>RW3ELDFMOFTES2EUAWQVCPPN7E</guid>
              </manager>
            </managers>
            <team_points>
              <coverage_type>week</coverage_type>
              <week>1</week>
              <total>103.82</total>
            </team_points>
            <team_projected_points>
              <coverage_type>week</coverage_type>
              <week>1</week>
              <total>110.41</total>
            </team_projected_points>
          </team>
        </teams>
      </matchup>
      <matchup>
        <week>5</week>
        <status>postevent</status>
        <is_tied>0</is_tied>
        <winner_team_key>223.l.431.t.1</winner_team_key>
        <teams count="2">
          <team>
            <team_key>223.l.431.t.1</team_key>
            <team_id>1</team_id>
            <name>PFW - Blunda</name>
            <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/1</url>
            <team_logos>
              <team_logo>
                <size>medium</size>
                <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_01_48.gif</url>
              </team_logo>
            </team_logos>
            <division_id>2</division_id>
            <faab_balance>22</faab_balance>
            <managers>
              <manager>
                <manager_id>13</manager_id>
                <nickname>Michael Blunda</nickname>
                <guid>XNAXQZRDZPJ3RVFMY7ZTSWEFLU</guid>
              </manager>
            </managers>
            <team_points>
              <coverage_type>week</coverage_type>
              <week>5</week>
              <total>140.00</total>
            </team_points>
            <team_projected_points>
              <coverage_type>week</coverage_type>
              <week>5</week>
              <total>110.85</total>
            </team_projected_points>
          </team>
          <team>
            <team_key>223.l.431.t.7</team_key>
            <team_id>7</team_id>
            <name>RotoWire_Liss</name>
            <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/7</url>
            <team_logos>
              <team_logo>
                <size>medium</size>
                <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_10_48.gif</url>
              </team_logo>
            </team_logos>
            <division_id>2</division_id>
            <faab_balance>68</faab_balance>
            <managers>
              <manager>
                <manager_id>7</manager_id>
                <nickname>RotoWire_Liss</nickname>
                <guid>4BDB5LIG3IFVROH7SRBX44LBZM</guid>
              </manager>
            </managers>
            <team_points>
              <coverage_type>week</coverage_type>
              <week>5</week>
              <total>86.47</total>
            </team_points>
            <team_projected_points>
              <coverage_type>week</coverage_type>
              <week>5</week>
              <total>88.14</total>
            </team_projected_points>
          </team>
        </teams>
      </matchup>
    </matchups>
  </team>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1/stats;type=season - team’s season stats in a NFL H2H league

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1/stats;type=season" time="129.66799736023ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <team>
    <team_key>223.l.431.t.1</team_key>
    <team_id>1</team_id>
    <name>PFW - Blunda</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431/1</url>
    <team_logos>
      <team_logo>
        <size>medium</size>
        <url>http://l.yimg.com/a/i/us/sp/fn/default/full/nfl/icon_01_48.gif</url>
      </team_logo>
    </team_logos>
    <division_id>2</division_id>
    <faab_balance>22</faab_balance>
    <managers>
      <manager>
        <manager_id>13</manager_id>
        <nickname>Michael Blunda</nickname>
        <guid>XNAXQZRDZPJ3RVFMY7ZTSWEFLU</guid>
      </manager>
    </managers>
    <team_points>
      <coverage_type>season</coverage_type>
      <season>2009</season>
      <total>1461.71</total>
    </team_points>
  </team>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/team/253.l.102614.t.10/stats;type=date;date=2011-07-06 - team’s date stats in a MLB roto league

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/team/253.l.102614.t.10/stats;date=2011-07-06;type=date" time="68.986892700195ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <team>
    <team_key>253.l.102614.t.10</team_key>
    <team_id>10</team_id>
    <name>Matt Dzaman</name>
    <url>http://baseball.fantasysports.yahoo.com/b1/102614/10</url>
    <team_logos>
      <team_logo>
        <size>medium</size>
        <url>http://l.yimg.com/a/i/us/sp/fn/mlb/gr/icon_12_2.gif</url>
      </team_logo>
    </team_logos>
    <managers>
      <manager>
        <manager_id>10</manager_id>
        <nickname>smock514</nickname>
        <guid>VZVEVUCLSJAHSM73FMJ4BYFIKU</guid>
        <exposed_yahoo_id>1</exposed_yahoo_id>
      </manager>
    </managers>
    <team_stats>
      <coverage_type>date</coverage_type>
      <date>2011-07-06</date>
      <stats>
        <stat>
          <stat_id>60</stat_id>
          <value>13/31</value>
        </stat>
        <stat>
          <stat_id>7</stat_id>
          <value>9</value>
        </stat>
        <stat>
          <stat_id>12</stat_id>
          <value>3</value>
        </stat>
        <stat>
          <stat_id>13</stat_id>
          <value>11</value>
        </stat>
        <stat>
          <stat_id>16</stat_id>
          <value>1</value>
        </stat>
        <stat>
          <stat_id>3</stat_id>
          <value>.419</value>
        </stat>
        <stat>
          <stat_id>50</stat_id>
          <value>7.0</value>
        </stat>
        <stat>
          <stat_id>28</stat_id>
          <value>1</value>
        </stat>
        <stat>
          <stat_id>32</stat_id>
          <value>0</value>
        </stat>
        <stat>
          <stat_id>42</stat_id>
          <value>6</value>
        </stat>
        <stat>
          <stat_id>26</stat_id>
          <value>1.29</value>
        </stat>
        <stat>
          <stat_id>27</stat_id>
          <value>0.71</value>
        </stat>
      </stats>
    </team_stats>
  </team>
</fantasy_content>
*/

/*
Roster resource¶

Description¶

Players on a team are organized into rosters corresponding to certain weeks, in NFL, or certain dates, in MLB, NBA, and NHL. Each player on a roster will be assigned a position if they’re in the starting lineup, or will be on the bench. You can only receive credit for stats accumulated by players in your starting lineup.

You can use this API to edit your lineup by PUTting up new positions for the players on a roster. You can also add/drop players from your roster by `POSTing new transactions <#transactions-collection-POST>`__ to the league’s transactions collection.

HTTP Operations Supported¶

GET
`PUT <#roster-resource-PUT>`__
URIs¶

http://fantasysports.yahooapis.com/fantasy/v2/team//roster

Any sub-resource under a roster is extracted using a URI like:

http://fantasysports.yahooapis.com/fantasy/v2/team//roster/

For NFL, you can specify a week parameter to retrieve a specific week – otherwise it will default to the current roster

http://fantasysports.yahooapis.com/fantasy/v2/team//roster;week=10

For MLB, NHL, or NBA, you can specify a date parameter to retrieve a specific date – otherwise it will default to today’s roster.

http://fantasysports.yahooapis.com/fantasy/v2/team//roster;date=2011-05-01

Sub-resources¶

Default sub-resource: players

Name	Description	URI	Sample
 	Access the players collection within the roster.	/fantasy/v2/team//roster/players	http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.9/roster/players
Sample XML¶

http://fantasysports.yahooapis.com/fantasy/v2/team/253.l.102614.t.10/roster/players

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/team/253.l.102614.t.10/roster/players" time="110.02206802368ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <team>
    <team_key>253.l.102614.t.10</team_key>
    <team_id>10</team_id>
    <name>Matt Dzaman</name>
    <url>http://baseball.fantasysports.yahoo.com/b1/102614/10</url>
    <team_logos>
      <team_logo>
        <size>medium</size>
        <url>http://l.yimg.com/a/i/us/sp/fn/mlb/gr/icon_12_2.gif</url>
      </team_logo>
    </team_logos>
    <managers>
      <manager>
        <manager_id>10</manager_id>
        <nickname>Sean Montgomery</nickname>
        <guid>VZVEVUCLSJAHSM73FMJ4BYFIKU</guid>
        <is_current_login>1</is_current_login>
      </manager>
    </managers>
    <roster>
      <coverage_type>date</coverage_type>
      <date>2011-07-22</date>
      <players count="22">
        <player>
          <player_key>253.p.7569</player_key>
          <player_id>7569</player_id>
          <name>
            <full>Brian McCann</full>
            <first>Brian</first>
            <last>McCann</last>
            <ascii_first>Brian</ascii_first>
            <ascii_last>McCann</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7569</editorial_player_key>
          <editorial_team_key>mlb.t.15</editorial_team_key>
          <editorial_team_full_name>Atlanta Braves</editorial_team_full_name>
          <editorial_team_abbr>Atl</editorial_team_abbr>
          <uniform_number>16</uniform_number>
          <display_position>C</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7569.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=eYxVIp_jg4DlEZmIgv6idg--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>C</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>C</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.7054</player_key>
          <player_id>7054</player_id>
          <name>
            <full>Adrian Gonzalez</full>
            <first>Adrian</first>
            <last>Gonzalez</last>
            <ascii_first>Adrian</ascii_first>
            <ascii_last>Gonzalez</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7054</editorial_player_key>
          <editorial_team_key>mlb.t.2</editorial_team_key>
          <editorial_team_full_name>Boston Red Sox</editorial_team_full_name>
          <editorial_team_abbr>Bos</editorial_team_abbr>
          <uniform_number>28</uniform_number>
          <display_position>1B</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7054.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=54BODgSe4P3NxShTjtIt9g--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>1B</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>1B</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.7746</player_key>
          <player_id>7746</player_id>
          <name>
            <full>Howie Kendrick</full>
            <first>Howie</first>
            <last>Kendrick</last>
            <ascii_first>Howie</ascii_first>
            <ascii_last>Kendrick</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7746</editorial_player_key>
          <editorial_team_key>mlb.t.3</editorial_team_key>
          <editorial_team_full_name>Los Angeles Angels</editorial_team_full_name>
          <editorial_team_abbr>LAA</editorial_team_abbr>
          <uniform_number>47</uniform_number>
          <display_position>1B,2B,OF</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7746.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=O01i1gfOs6RgisJQjmdipQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>1B</position>
            <position>2B</position>
            <position>OF</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <has_recent_player_notes>1</has_recent_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>2B</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>0</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.7737</player_key>
          <player_id>7737</player_id>
          <name>
            <full>Martin Prado</full>
            <first>Martin</first>
            <last>Prado</last>
            <ascii_first>Martin</ascii_first>
            <ascii_last>Prado</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7737</editorial_player_key>
          <editorial_team_key>mlb.t.15</editorial_team_key>
          <editorial_team_full_name>Atlanta Braves</editorial_team_full_name>
          <editorial_team_abbr>Atl</editorial_team_abbr>
          <uniform_number>14</uniform_number>
          <display_position>2B,3B,OF</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7737.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=WPYI1xO62JwsL8QturlmJw--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>2B</position>
            <position>3B</position>
            <position>OF</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>3B</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.7744</player_key>
          <player_id>7744</player_id>
          <name>
            <full>Erick Aybar</full>
            <first>Erick</first>
            <last>Aybar</last>
            <ascii_first>Erick</ascii_first>
            <ascii_last>Aybar</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7744</editorial_player_key>
          <editorial_team_key>mlb.t.3</editorial_team_key>
          <editorial_team_full_name>Los Angeles Angels</editorial_team_full_name>
          <editorial_team_abbr>LAA</editorial_team_abbr>
          <uniform_number>2</uniform_number>
          <display_position>SS</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7744.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=qHzsNyGFtGYxlpMxtysSPQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>SS</position>
            <position>Util</position>
          </eligible_positions>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>SS</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.7977</player_key>
          <player_id>7977</player_id>
          <name>
            <full>Andrew McCutchen</full>
            <first>Andrew</first>
            <last>McCutchen</last>
            <ascii_first>Andrew</ascii_first>
            <ascii_last>McCutchen</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7977</editorial_player_key>
          <editorial_team_key>mlb.t.23</editorial_team_key>
          <editorial_team_full_name>Pittsburgh Pirates</editorial_team_full_name>
          <editorial_team_abbr>Pit</editorial_team_abbr>
          <uniform_number>22</uniform_number>
          <display_position>OF</display_position>
          <image_url>http://l.yimg.com/a/p/sp/tools/med/2011/05/ipt/1304541420.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=61GeaeZwqXZWy2ITOX62Zg--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>OF</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>OF</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.7104</player_key>
          <player_id>7104</player_id>
          <name>
            <full>Shane Victorino</full>
            <first>Shane</first>
            <last>Victorino</last>
            <ascii_first>Shane</ascii_first>
            <ascii_last>Victorino</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7104</editorial_player_key>
          <editorial_team_key>mlb.t.22</editorial_team_key>
          <editorial_team_full_name>Philadelphia Phillies</editorial_team_full_name>
          <editorial_team_abbr>Phi</editorial_team_abbr>
          <uniform_number>8</uniform_number>
          <display_position>OF</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7104.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=QE9iNVRK5VCHq650WQii4g--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>OF</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>OF</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.8239</player_key>
          <player_id>8239</player_id>
          <name>
            <full>Matt Joyce</full>
            <first>Matt</first>
            <last>Joyce</last>
            <ascii_first>Matt</ascii_first>
            <ascii_last>Joyce</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8239</editorial_player_key>
          <editorial_team_key>mlb.t.30</editorial_team_key>
          <editorial_team_full_name>Tampa Bay Rays</editorial_team_full_name>
          <editorial_team_abbr>TB</editorial_team_abbr>
          <uniform_number>20</uniform_number>
          <display_position>OF</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8239.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=ZIy1Z9IryxkYXVsSsodRfQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>OF</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>OF</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.8857</player_key>
          <player_id>8857</player_id>
          <name>
            <full>Eric Hosmer</full>
            <first>Eric</first>
            <last>Hosmer</last>
            <ascii_first>Eric</ascii_first>
            <ascii_last>Hosmer</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8857</editorial_player_key>
          <editorial_team_key>mlb.t.7</editorial_team_key>
          <editorial_team_full_name>Kansas City Royals</editorial_team_full_name>
          <editorial_team_abbr>KC</editorial_team_abbr>
          <uniform_number>35</uniform_number>
          <display_position>1B</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110706/8857.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=h5CchQfLoumJ6xXRYqQNOw--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>1B</position>
            <position>Util</position>
          </eligible_positions>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>Util</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.8171</player_key>
          <player_id>8171</player_id>
          <name>
            <full>Jay Bruce</full>
            <first>Jay</first>
            <last>Bruce</last>
            <ascii_first>Jay</ascii_first>
            <ascii_last>Bruce</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8171</editorial_player_key>
          <editorial_team_key>mlb.t.17</editorial_team_key>
          <editorial_team_full_name>Cincinnati Reds</editorial_team_full_name>
          <editorial_team_abbr>Cin</editorial_team_abbr>
          <uniform_number>32</uniform_number>
          <display_position>OF</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8171.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=icy.cvuP8XXvyrQKm7m3HA--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>OF</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <has_recent_player_notes>1</has_recent_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>BN</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>0</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.8401</player_key>
          <player_id>8401</player_id>
          <name>
            <full>Elvis Andrus</full>
            <first>Elvis</first>
            <last>Andrus</last>
            <ascii_first>Elvis</ascii_first>
            <ascii_last>Andrus</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8401</editorial_player_key>
          <editorial_team_key>mlb.t.13</editorial_team_key>
          <editorial_team_full_name>Texas Rangers</editorial_team_full_name>
          <editorial_team_abbr>Tex</editorial_team_abbr>
          <uniform_number>1</uniform_number>
          <display_position>SS</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8401.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=HIAp3xabHCwOw.hJpkbd1w--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>SS</position>
            <position>Util</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <has_recent_player_notes>1</has_recent_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>BN</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.7926</player_key>
          <player_id>7926</player_id>
          <name>
            <full>Yovani Gallardo</full>
            <first>Yovani</first>
            <last>Gallardo</last>
            <ascii_first>Yovani</ascii_first>
            <ascii_last>Gallardo</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7926</editorial_player_key>
          <editorial_team_key>mlb.t.8</editorial_team_key>
          <editorial_team_full_name>Milwaukee Brewers</editorial_team_full_name>
          <editorial_team_abbr>Mil</editorial_team_abbr>
          <uniform_number>49</uniform_number>
          <display_position>SP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7926.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=1lRXgDptEQng1WvYIB3vDQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>SP</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.7172</player_key>
          <player_id>7172</player_id>
          <name>
            <full>Dan Haren</full>
            <first>Dan</first>
            <last>Haren</last>
            <ascii_first>Dan</ascii_first>
            <ascii_last>Haren</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7172</editorial_player_key>
          <editorial_team_key>mlb.t.3</editorial_team_key>
          <editorial_team_full_name>Los Angeles Angels</editorial_team_full_name>
          <editorial_team_abbr>LAA</editorial_team_abbr>
          <uniform_number>24</uniform_number>
          <display_position>SP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7172.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=fnc4Dr.qGpHVMT8tW4phOQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>SP</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.6210</player_key>
          <player_id>6210</player_id>
          <name>
            <full>Kyle Farnsworth</full>
            <first>Kyle</first>
            <last>Farnsworth</last>
            <ascii_first>Kyle</ascii_first>
            <ascii_last>Farnsworth</ascii_last>
          </name>
          <editorial_player_key>mlb.p.6210</editorial_player_key>
          <editorial_team_key>mlb.t.30</editorial_team_key>
          <editorial_team_full_name>Tampa Bay Rays</editorial_team_full_name>
          <editorial_team_abbr>TB</editorial_team_abbr>
          <uniform_number>43</uniform_number>
          <display_position>RP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/6210.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=kJYLFUywffdzSTtTQ5MupQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>RP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>RP</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.8929</player_key>
          <player_id>8929</player_id>
          <name>
            <full>Javy Guerra</full>
            <first>Javy</first>
            <last>Guerra</last>
            <ascii_first>Javy</ascii_first>
            <ascii_last>Guerra</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8929</editorial_player_key>
          <editorial_team_key>mlb.t.19</editorial_team_key>
          <editorial_team_full_name>Los Angeles Dodgers</editorial_team_full_name>
          <editorial_team_abbr>LAD</editorial_team_abbr>
          <uniform_number>54</uniform_number>
          <display_position>RP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/blank_player2.gif?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=8G0MjQyD1AdYbnv.fd2Wog--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>RP</position>
            <position>P</position>
          </eligible_positions>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>RP</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.7279</player_key>
          <player_id>7279</player_id>
          <name>
            <full>Jesse Crain</full>
            <first>Jesse</first>
            <last>Crain</last>
            <ascii_first>Jesse</ascii_first>
            <ascii_last>Crain</ascii_last>
          </name>
          <editorial_player_key>mlb.p.7279</editorial_player_key>
          <editorial_team_key>mlb.t.4</editorial_team_key>
          <editorial_team_full_name>Chicago White Sox</editorial_team_full_name>
          <editorial_team_abbr>CWS</editorial_team_abbr>
          <uniform_number>26</uniform_number>
          <display_position>RP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7279.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=Xx9Emr_lBK3smmABr7fOcg--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>RP</position>
            <position>P</position>
          </eligible_positions>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>P</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.8193</player_key>
          <player_id>8193</player_id>
          <name>
            <full>Max Scherzer</full>
            <first>Max</first>
            <last>Scherzer</last>
            <ascii_first>Max</ascii_first>
            <ascii_last>Scherzer</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8193</editorial_player_key>
          <editorial_team_key>mlb.t.6</editorial_team_key>
          <editorial_team_full_name>Detroit Tigers</editorial_team_full_name>
          <editorial_team_abbr>Det</editorial_team_abbr>
          <uniform_number>37</uniform_number>
          <display_position>SP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8193.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=lcgNVFPn0gpchY5fCbb6nw--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>P</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
        <player>
          <player_key>253.p.8099</player_key>
          <player_id>8099</player_id>
          <name>
            <full>Ian Kennedy</full>
            <first>Ian</first>
            <last>Kennedy</last>
            <ascii_first>Ian</ascii_first>
            <ascii_last>Kennedy</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8099</editorial_player_key>
          <editorial_team_key>mlb.t.29</editorial_team_key>
          <editorial_team_full_name>Arizona Diamondbacks</editorial_team_full_name>
          <editorial_team_abbr>Ari</editorial_team_abbr>
          <uniform_number>31</uniform_number>
          <display_position>SP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8099.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=TjXYM8e9wtrcLfrhDnCMKQ--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <has_recent_player_notes>1</has_recent_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>P</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.8179</player_key>
          <player_id>8179</player_id>
          <name>
            <full>Gio Gonzalez</full>
            <first>Gio</first>
            <last>Gonzalez</last>
            <ascii_first>Gio</ascii_first>
            <ascii_last>Gonzalez</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8179</editorial_player_key>
          <editorial_team_key>mlb.t.11</editorial_team_key>
          <editorial_team_full_name>Oakland Athletics</editorial_team_full_name>
          <editorial_team_abbr>Oak</editorial_team_abbr>
          <uniform_number>47</uniform_number>
          <display_position>SP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8179.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=Wg7KqIjVG4zwo0znBbEViw--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>BN</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.8759</player_key>
          <player_id>8759</player_id>
          <name>
            <full>Michael Pineda</full>
            <first>Michael</first>
            <last>Pineda</last>
            <ascii_first>Michael</ascii_first>
            <ascii_last>Pineda</ascii_last>
          </name>
          <editorial_player_key>mlb.p.8759</editorial_player_key>
          <editorial_team_key>mlb.t.12</editorial_team_key>
          <editorial_team_full_name>Seattle Mariners</editorial_team_full_name>
          <editorial_team_abbr>Sea</editorial_team_abbr>
          <uniform_number>36</uniform_number>
          <display_position>SP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/8759.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=31SmIDWcet4v3AVAOGrY2g--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>P</position>
          </eligible_positions>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>BN</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.6571</player_key>
          <player_id>6571</player_id>
          <name>
            <full>Ryan Vogelsong</full>
            <first>Ryan</first>
            <last>Vogelsong</last>
            <ascii_first>Ryan</ascii_first>
            <ascii_last>Vogelsong</ascii_last>
          </name>
          <editorial_player_key>mlb.p.6571</editorial_player_key>
          <editorial_team_key>mlb.t.26</editorial_team_key>
          <editorial_team_full_name>San Francisco Giants</editorial_team_full_name>
          <editorial_team_abbr>SF</editorial_team_abbr>
          <uniform_number>32</uniform_number>
          <display_position>SP,RP</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110706/6571.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=bdeeFeFntdasbz_0xzXCGA--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>P</position_type>
          <eligible_positions>
            <position>SP</position>
            <position>RP</position>
            <position>P</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>BN</position>
          </selected_position>
        </player>
        <player>
          <player_key>253.p.7382</player_key>
          <player_id>7382</player_id>
          <name>
            <full>David Wright</full>
            <first>David</first>
            <last>Wright</last>
            <ascii_first>David</ascii_first>
            <ascii_last>Wright</ascii_last>
          </name>
          <status>DL</status>
          <on_disabled_list>1</on_disabled_list>
          <editorial_player_key>mlb.p.7382</editorial_player_key>
          <editorial_team_key>mlb.t.21</editorial_team_key>
          <editorial_team_full_name>New York Mets</editorial_team_full_name>
          <editorial_team_abbr>NYM</editorial_team_abbr>
          <uniform_number>5</uniform_number>
          <display_position>3B</display_position>
          <image_url>http://l.yimg.com/a/i/us/sp/v/mlb/players_l/20110503x/7382.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=QNOFMSgR6NuPxUwDMUSM1w--</image_url>
          <is_undroppable>0</is_undroppable>
          <position_type>B</position_type>
          <eligible_positions>
            <position>3B</position>
            <position>Util</position>
            <position>DL</position>
          </eligible_positions>
          <has_player_notes>1</has_player_notes>
          <has_recent_player_notes>1</has_recent_player_notes>
          <selected_position>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <position>DL</position>
          </selected_position>
          <starting_status>
            <coverage_type>date</coverage_type>
            <date>2011-07-22</date>
            <is_starting>1</is_starting>
          </starting_status>
        </player>
      </players>
    </roster>
  </team>
</fantasy_content>
PUT¶

Using PUT, you may modify a subset of players on the roster for a particular day, specifically in terms of changing their position or whether they’re in the starting lineup. The URL for PUTting to a Roster resource is:

http://fantasysports.yahooapis.com/fantasy/v2/team//roster

You may move as many players as you like in your input XML – any players whose position you do not change will stay in the same position they were previously. If you try to move players in an invalid way, you will receive an error and no changes will be made.

Your input XML should look like:

NFL:

<?xml version="1.0"?>
<fantasy_content>
  <roster>
    <coverage_type>week</coverage_type>
    <week>13</week>

    <players>
      <player>
        <player_key>242.p.8332</player_key>
        <position>WR</position>
      </player>
      <player>
        <player_key>242.p.1423</player_key>
        <position>BN</position>
      </player>
    </players>
  </roster>
</fantasy_content>
MLB, NBA, or NHL:

<?xml version="1.0"?>
<fantasy_content>
  <roster>
    <coverage_type>date</coverage_type>
    <date>2011-05-01</date>

    <players>
      <player>
        <player_key>253.p.8332</player_key>
        <position>1B</position>
      </player>
      <player>
        <player_key>253.p.1423</player_key>
        <position>BN</position>
      </player>
    </players>
  </roster>
</fantasy_content>
*/

/*
Teams collection¶

Description¶

With the Teams API, you can obtain information from a collection of teams simultaneously. The teams collection is qualified in the URI by a particular league to obtain information about teams within the league, or by a particular user (and optionally, a game) to obtain information about the teams owned by the user. Each element beneath the Teams Collection will be a Team Resource

HTTP Operations Supported¶

GET
URIs¶

URI	Description	Sample
http://fantasysports.yahooapis.com/fantasy/v2/league//teams	Fetch all teams within a league.	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/teams
http://fantasysports.yahooapis.com/fantasy/v2/teams;team_keys=,{team_key2}	Fetch specific teams {team_key1} and {team_key2}	http://fantasysports.yahooapis.com/fantasy/v2/teams;team_keys=223.l.431.t.1,223.l.431.t.2
http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys={league_key1},{league_key2}/teams	Fetch all teams of the leagues {league_key1} and {league_key2}	http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=223.l.431,223.l.21821/teams
http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/teams	Fetch all teams for the logged in user	http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/teams
http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys={game_key1},{game_key2}/teams	Fetch all teams for the logged in user for the games {game_key1} and {game_key2}	http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys=nfl,mlb/teams
Any sub-resource valid for a team is a valid sub-resource under the teams collection.

Any sub-resource for a collection of teams is extracted using a URI like:

/teams/{sub_resource}

OR

/teams;team_keys={team_key1},{team_key2}/{sub_resource}

Multiple sub-resources can be extracted from teams in the same URI using a format like:

/teams;out={sub_resource_1},{sub_resource_2}

OR

/teams;team_keys={team_key1},{team_key2};out={sub_resource_1},{sub_resource_2}
*/

/*
Player resource¶

Description¶

With the Player API, you can obtain the player (athlete) related information, such as their name, professional team, and eligible positions. The player is identified in the context of a particular game, and can be requested as the base of your URI by using the global ````.

HTTP Operations Supported¶

GET
URIs¶

http://fantasysports.yahooapis.com/fantasy/v2/player/

Any sub-resource under a player is extracted using a URI like:

http://fantasysports.yahooapis.com/fantasy/v2/player//

Multiple sub-resources can be extracted from player in the same URI using a format like:

http://fantasysports.yahooapis.com/fantasy/v2/player/;out=,{sub_resource_2}

Player key format¶

.p.{player_id}

Example:pnfl.p.5479 or 223.p.5479

Sub-resources¶

Default sub-resource: metadata

Name	Description	URI	Sample
metadata	Includes player key, id, name, editorial information, image, eligible positions, etc.	/fantasy/v2/player//metadata	Drew Brees’s info in the 2009 season: http://fantasysports.yahooapis.com/fantasy/v2/player/223.p.5479
stats	Player stats and points (if in a league context).
Season stats:/fantasy/v2/player//stats

Week stats: /fantasy/v2/player//stats;type=week;week={week}

Here {week} is a non-zero integer.

Drew Brees’s info and stats in the 2009 season: http://fantasysports.yahooapis.com/fantasy/v2/player/223.p.5479/stats
ownership	The player ownership status within a league (whether they’re owned by a team, on waivers, or free agents). Only relevant within a league.	/fantasy/v2/league//players;player_keys=/ownership	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players;player_keys=223.p.5479/ownership
percent_owned	Data about ownership percentage of the player	/fantasy/v2/player//percent_owned	The percentage of leagues in which Drew Brees was owned in the 2009 game: http://fantasysports.yahooapis.com/fantasy/v2/player/223.p.5479/percent_owned
draft_analysis	Average pick, Average round and Percent Drafted.	/fantasy/v2/player//draft_analysis	Yahoo! fantasy draft information for Drew Brees in 2009: http://fantasysports.yahooapis.com/fantasy/v2/player/223.p.5479/draft_analysis
Sample XML¶

http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players;player_keys=223.p.5479 - Player in a NFL league context

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players;player_keys=223.p.5479" time="508.72206687927ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <league>
    <league_key>223.l.431</league_key>
    <league_id>431</league_id>
    <name>Y! Friends and Family League</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431</url>
    <password>liss</password>
    <draft_status>postdraft</draft_status>
    <num_teams>14</num_teams>
    <edit_key>17</edit_key>
    <weekly_deadline/>
    <league_update_timestamp>1262595518</league_update_timestamp>
    <scoring_type>head</scoring_type>
    <current_week>16</current_week>
    <start_week>1</start_week>
    <end_week>16</end_week>
    <is_finished>1</is_finished>
    <players count="1">
      <player>
        <player_key>223.p.5479</player_key>
        <player_id>5479</player_id>
        <name>
          <full>Drew Brees</full>
          <first>Drew</first>
          <last>Brees</last>
          <ascii_first>Drew</ascii_first>
          <ascii_last>Brees</ascii_last>
        </name>
        <status>P</status>
        <editorial_player_key>nfl.p.5479</editorial_player_key>
        <editorial_team_key>nfl.t.18</editorial_team_key>
        <editorial_team_full_name>New Orleans Saints</editorial_team_full_name>
        <editorial_team_abbr>NO</editorial_team_abbr>
        <bye_weeks>
          <week>5</week>
        </bye_weeks>
        <uniform_number>9</uniform_number>
        <display_position>QB</display_position>
        <image_url>http://l.yimg.com/a/i/us/sp/v/nfl/players_l/headshots/20100903/5479.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=LTUFmLVwQ.kvKzbhSsG94w--</image_url>
        <is_undroppable>0</is_undroppable>
        <position_type>O</position_type>
        <eligible_positions>
          <position>QB</position>
        </eligible_positions>
      </player>
    </players>
  </league>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players;player_keys=223.p.5479/stats - Player season stats in a NFL league context

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players;player_keys=223.p.5479/stats" time="3140.1500701904ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <league>
    <league_key>223.l.431</league_key>
    <league_id>431</league_id>
    <name>Y! Friends and Family League</name>
    <url>http://football.fantasysports.yahoo.com/archive/pnfl/2009/431</url>
    <password>liss</password>
    <draft_status>postdraft</draft_status>
    <num_teams>14</num_teams>
    <edit_key>17</edit_key>
    <weekly_deadline/>
    <league_update_timestamp>1262595518</league_update_timestamp>
    <scoring_type>head</scoring_type>
    <current_week>16</current_week>
    <start_week>1</start_week>
    <end_week>16</end_week>
    <is_finished>1</is_finished>
    <players count="1">
      <player>
        <player_key>223.p.5479</player_key>
        <player_id>5479</player_id>
        <name>
          <full>Drew Brees</full>
          <first>Drew</first>
          <last>Brees</last>
          <ascii_first>Drew</ascii_first>
          <ascii_last>Brees</ascii_last>
        </name>
        <status>P</status>
        <editorial_player_key>nfl.p.5479</editorial_player_key>
        <editorial_team_key>nfl.t.18</editorial_team_key>
        <editorial_team_full_name>New Orleans Saints</editorial_team_full_name>
        <editorial_team_abbr>NO</editorial_team_abbr>
        <bye_weeks>
          <week>5</week>
        </bye_weeks>
        <uniform_number>9</uniform_number>
        <display_position>QB</display_position>
        <image_url>http://l.yimg.com/a/i/us/sp/v/nfl/players_l/headshots/20100903/5479.jpg?x=46&amp;y=60&amp;xc=1&amp;yc=1&amp;wc=164&amp;hc=215&amp;q=100&amp;sig=LTUFmLVwQ.kvKzbhSsG94w--</image_url>
        <is_undroppable>0</is_undroppable>
        <position_type>O</position_type>
        <eligible_positions>
          <position>QB</position>
        </eligible_positions>
        <player_stats>
          <coverage_type>season</coverage_type>
          <season>2009</season>
          <stats>
            <stat>
              <stat_id>4</stat_id>
              <value>4388</value>
            </stat>
            <stat>
              <stat_id>5</stat_id>
              <value>34</value>
            </stat>
            <stat>
              <stat_id>6</stat_id>
              <value>11</value>
            </stat>
            <stat>
              <stat_id>9</stat_id>
              <value>33</value>
            </stat>
            <stat>
              <stat_id>10</stat_id>
              <value>2</value>
            </stat>
            <stat>
              <stat_id>11</stat_id>
              <value>1</value>
            </stat>
            <stat>
              <stat_id>12</stat_id>
              <value>-4</value>
            </stat>
            <stat>
              <stat_id>13</stat_id>
              <value>0</value>
            </stat>
            <stat>
              <stat_id>15</stat_id>
              <value>0</value>
            </stat>
            <stat>
              <stat_id>16</stat_id>
              <value>0</value>
            </stat>
            <stat>
              <stat_id>18</stat_id>
              <value>6</value>
            </stat>
            <stat>
              <stat_id>57</stat_id>
              <value>0</value>
            </stat>
          </stats>
        </player_stats>
        <player_points>
          <coverage_type>season</coverage_type>
          <season>2009</season>
          <total>310.17</total>
        </player_points>
      </player>
    </players>
  </league>
</fantasy_content>
*/

/*
Players collection¶

Description¶

With the Players API, you can obtain information from a collection of players simultaneously. To obtains general players information, the players collection can be qualified in the URI by a particular game, league or team. To obtain specific league or team related information, the players collection is qualified by the relevant league or team. Each element beneath the Players Collection will be a Player Resource

HTTP Operations Supported¶

GET
URIs¶

URI	Description	Sample
http://fantasysports.yahooapis.com/fantasy/v2/league//players	Fetch all players within a league.	http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/players
http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys={league_key1},{league_key2}/players	Fetch all players from the leagues {league_key1} and {league_key2}	http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=223.l.431,223.l.21821/players
http://fantasysports.yahooapis.com/fantasy/v2/team//players	Fetch all players within a team.	http://fantasysports.yahooapis.com/fantasy/v2/team/223.l.431.t.1/players
http://fantasysports.yahooapis.com/fantasy/v2/teams;team_keys={team_key1},{team_key2}/players	Fetch all players from the teams {team_key1} and {team_key2}	http://fantasysports.yahooapis.com/fantasy/v2/teams;team_keys=223.l.431.t.1,223.l.431.t.2/players
http://fantasysports.yahooapis.com/fantasy/v2/players;player_keys=,{player_key2}	Fetch specific players {player_key1} and {player_key2}	http://fantasysports.yahooapis.com/fantasy/v2/players;player_keys=223.p.5479,223.p.1025
Any sub-resource valid for a player is a valid sub-resource under the players collection.

Any sub-resource for a collection of players is extracted using a URI like:

/players/{sub_resource}

OR

/players;player_keys={player_key1},{player_key2}/{sub_resource}

Multiple sub-resources can be extracted from players in the same URI using a format like:

/players;out={sub_resource_1},{sub_resource_2}

OR

/players;player_keys={player_key1},{player_key2};out={sub_resource_1},{sub_resource_2}

Filters¶

The players collection can have filters such as the following to obtain a subset of a players collection that satisfy the filtering condition. The filters can be combined to obtain a more restricted list of players.

Filter parameter	Filter parameter values	Usage
position	Valid player positions
/players;position=QB
Note

Applied only in a league’s context

status
A (all available players)

FA (free agents only)

W (waivers only)

T (all taken players)

K (keepers only)

/players;status=A
Note

Applied only in a league’s context

search	player name
/players;search=smith
Note

Applied only in a league’s context

sort
{stat_id}

NAME (last, first)

OR (overall rank)

AR (actual rank)

PTS (fantasy points)

/players;sort=60
Note

Applied only in a league’s context

sort_type
season

date (baseball, basketball, and hockey only)

week (football only)

lastweek (baseball, basketball, and hockey only)

lastmonth

/players;sort_type=season
Note

Applied only in a league’s context

sort_season	year
/players;sort_type=season;sort_season=2010
Note

Applied only in a league’s context

sort_date (baseball, basketball, and hockey only)	YYYY-MM-DD
/players;sort_type=date;sort_date=2010-02-01
Note

Applied only in a league’s context

sort_week (football only)	week
/players;sort_type=week;sort_week=10
Note

Applied only in a league’s context

start	Any integer 0 or greater	/players;start=25
count	Any integer greater than 0	/players;count=5
*/

/*
Transaction resource¶

Description¶

With the Transaction API, you can obtain information about transactions (adds, drops, trades, and league settings changes) performed on a league. A transaction is identified in the context of a particular league, although you can request a particular Transaction Resource as the base of your URI by using the global ````.

You can also PUT to the API to perform operations like editing waiver priorities or FAAB bids, or modifying the state of pending trades. You can also cancel pending transactions by DELETEing them.

Keep in mind, if you don’t have the ```` for a waiver claim or pending trade, the only way to discover these transactions is to filter the league Transactions collection by a particular type (waiver or pending_trade) and by a particular ````. Pending transactions will not show up if you simply ask for all of the transactions in the league, because they can only be seen by certain teams.

HTTP Operations Supported¶

GET
`PUT <#transaction-resource-PUT>`__
`DELETE <#transaction-resource-DELETE>`__
URIs¶

http://fantasysports.yahooapis.com/fantasy/v2/transaction/

Any sub-resource under a transaction is extracted using a URI like:

http://fantasysports.yahooapis.com/fantasy/v2/transaction//

Multiple sub-resources can be extracted from transaction in the same URI using a format like:

http://fantasysports.yahooapis.com/fantasy/v2/transaction/;out=,{sub_resource_2}

Transaction key format¶

Completed transactions: .l.{league_id}.tr.{transaction_id}

Example:pnfl.l.431.tr.26 or 223.l.431.tr.26

Waiver claims: .l.{league_id}.w.c.{claim_id}

Example:257.l.193.w.c.2_6390

Pending trades: .l.{league_id}.pt.{pending_trade_id}

Example:257.l.193.pt.1

Sub-resources¶

Default sub-resources: metadata, players

Name	Description	URI	Sample
metadata	Includes transaction key, id, type, timestamp, status, players (not displayed for all transaction types)	/fantasy/v2/transaction//metadata	An add/drop transaction: http://fantasysports.yahooapis.com/fantasy/v2/transaction/223.l.431.tr.26
``

``
Players that are part of the transaction. The Player Resources will include a transaction data element by default.	/fantasy/v2/transaction//players	http://fantasysports.yahooapis.com/fantasy/v2/transaction/223.l.431.tr.26/players
Sample XML¶

http://fantasysports.yahooapis.com/fantasy/v2/transaction/257.l.193.tr.2 - Completed add/drop transaction

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/transaction/257.l.193.tr.2" time="51.784038543701ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <transaction>
    <transaction_key>257.l.193.tr.2</transaction_key>
    <transaction_id>2</transaction_id>
    <type>add/drop</type>
    <status>successful</status>
    <timestamp>1310694660</timestamp>
    <players count="2">
      <player>
        <player_key>257.p.7847</player_key>
        <player_id>7847</player_id>
        <name>
          <full>Owen Daniels</full>
          <first>Owen</first>
          <last>Daniels</last>
          <ascii_first>Owen</ascii_first>
          <ascii_last>Daniels</ascii_last>
        </name>
        <transaction_data>
          <type>add</type>
          <source_type>freeagents</source_type>
          <destination_type>team</destination_type>
          <destination_team_key>257.l.193.t.1</destination_team_key>
        </transaction_data>
      </player>
      <player>
        <player_key>257.p.6390</player_key>
        <player_id>6390</player_id>
        <name>
          <full>Anquan Boldin</full>
          <first>Anquan</first>
          <last>Boldin</last>
          <ascii_first>Anquan</ascii_first>
          <ascii_last>Boldin</ascii_last>
        </name>
        <transaction_data>
          <type>drop</type>
          <source_type>team</source_type>
          <source_team_key>257.l.193.t.1</source_team_key>
          <destination_type>waivers</destination_type>
        </transaction_data>
      </player>
    </players>
  </transaction>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/transaction/257.l.193.w.c.2_6390 - Waiver claim transaction

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/transaction/257.l.193.w.c.2_6390" time="30.953884124756ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <transaction>
    <transaction_key>257.l.193.w.c.2_6390</transaction_key>
    <type>waiver</type>
    <status>pending</status>
    <waiver_player_key>257.p.6390</waiver_player_key>
    <waiver_team_key>257.l.193.t.2</waiver_team_key>
    <waiver_date>2011-07-17</waiver_date>
    <waiver_priority>1</waiver_priority>
    <players count="2">
      <player>
        <player_key>257.p.6390</player_key>
        <player_id>6390</player_id>
        <name>
          <full>Anquan Boldin</full>
          <first>Anquan</first>
          <last>Boldin</last>
          <ascii_first>Anquan</ascii_first>
          <ascii_last>Boldin</ascii_last>
        </name>
        <transaction_data>
          <type>add</type>
          <source_type>waivers</source_type>
          <destination_type>team</destination_type>
          <destination_team_key>257.l.193.t.2</destination_team_key>
        </transaction_data>
      </player>
      <player>
        <player_key>257.p.8266</player_key>
        <player_id>8266</player_id>
        <name>
          <full>Marshawn Lynch</full>
          <first>Marshawn</first>
          <last>Lynch</last>
          <ascii_first>Marshawn</ascii_first>
          <ascii_last>Lynch</ascii_last>
        </name>
        <transaction_data>
          <type>drop</type>
          <source_type>team</source_type>
          <source_team_key>257.l.193.t.2</source_team_key>
          <destination_type>waivers</destination_type>
        </transaction_data>
      </player>
    </players>
  </transaction>
</fantasy_content>
http://fantasysports.yahooapis.com/fantasy/v2/transaction/257.l.193.pt.1 - Pending trade transaction

<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng" xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/transaction/257.l.193.pt.1" time="45.558929443359ms" copyright="Data provided by Yahoo! and STATS, LLC">
  <transaction>
    <transaction_key>257.l.193.pt.1</transaction_key>
    <type>pending_trade</type>
    <status>proposed</status>
    <trader_team_key>257.l.193.t.2</trader_team_key>
    <tradee_team_key>257.l.193.t.1</tradee_team_key>
    <trade_proposed_time>1310694832</trade_proposed_time>
    <trade_note>This is a great trade, fo' shizzle.</trade_note>
    <players count="2">
      <player>
        <player_key>257.p.8261</player_key>
        <player_id>8261</player_id>
        <name>
          <full>Adrian Peterson</full>
          <first>Adrian</first>
          <last>Peterson</last>
          <ascii_first>Adrian</ascii_first>
          <ascii_last>Peterson</ascii_last>
        </name>
        <transaction_data>
          <type>pending_trade</type>
          <source_type>team</source_type>
          <source_team_key>257.l.193.t.2</source_team_key>
          <destination_type>team</destination_type>
          <destination_team_key>257.l.193.t.1</destination_team_key>
        </transaction_data>
      </player>
      <player>
        <player_key>257.p.9527</player_key>
        <player_id>9527</player_id>
        <name>
          <full>Arian Foster</full>
          <first>Arian</first>
          <last>Foster</last>
          <ascii_first>Arian</ascii_first>
          <ascii_last>Foster</ascii_last>
        </name>
        <transaction_data>
          <type>pending_trade</type>
          <source_type>team</source_type>
          <source_team_key>257.l.193.t.1</source_team_key>
          <destination_type>team</destination_type>
          <destination_team_key>257.l.193.t.2</destination_team_key>
        </transaction_data>
      </player>
    </players>
  </transaction>
</fantasy_content>
*/

// PUT
// Using PUT, you may edit the waiver priority or FAAB bid for any of your
// pending waiver claims. You can also accept or reject trades that have been
// proposed to you, and allow or vote against trades if your league settings
// allow it. The URL for PUTting to a Transaction resource is:
// http://fantasysports.yahooapis.com/fantasy/v2/transaction/
//
// You can only PUT to Transactions of the types waiver or pending_trade.
//
// Editing Waivers
// Once you have the transaction_key for a waiver claim, which you can get by
// asking the transactions collection for all waivers for a certain team, you
// can edit the waiver priority or FAAB bid. The input XML should look like:
//     <?xml version='1.0'?>
//     <fantasy_content>
//       <transaction>
//         <transaction_key>248.l.55438.w.c.2_6093</transaction_key>
//         <type>waiver</type>
//         <waiver_priority>1</waiver_priority>
//         <faab_bid>20</faab_bid>
//       </transaction>
//     </fantasy_content>
func (y *YahooConfig) EditWaivers() {
	// PUT
}

// Accepting Trades
// Once you have the transaction_key for a pending trade that has been proposed
// to you, which you can get by asking the transactions collection for all
// pending trades for your team, you can choose to accept it. The input XML
// should look like:
//     <?xml version='1.0'?>
//     <fantasy_content>
//       <transaction>
//         <transaction_key>248.l.55438.pt.11</transaction_key>
//         <type>pending_trade</type>
//         <action>accept</action>
//         <trade_note>Dude, that is a totally fair trade.</trade_note>
//       </transaction>
//     </fantasy_content>
func (y *YahooConfig) AcceptTrade() {
	// PUT
}

// Rejecting Trades
// To reject a pending trade proposed to you, the input XML should look like:
//     <?xml version='1.0'?>
//     <fantasy_content>
//       <transaction>
//         <transaction_key>248.l.55438.pt.11</transaction_key>
//         <type>pending_trade</type>
//         <action>reject</action>
//         <trade_note>No way!</trade_note>
//       </transaction>
//     </fantasy_content>
func (y *YahooConfig) RejectTrade() {
	// PUT
}

// Allowing/Disallowing Trades
// If there are accepted trades in your league waiting to be processed, which
// you can get by asking the transactions collection for all pending trades for
// your team, and you’re the commissioner of a league that has the commissioner
// approve trades, you can choose to allow or disallow the trade. The input XML
// should look like:
//    <?xml version='1.0'?>
//    <fantasy_content>
//      <transaction>
//        <transaction_key>248.l.55438.pt.11</transaction_key>
//        <type>pending_trade</type>
//        <action>allow</action>
//      </transaction>
//    </fantasy_content>
// Or
//    <?xml version='1.0'?>
//    <fantasy_content>
//      <transaction>
//        <transaction_key>248.l.55438.pt.11</transaction_key>
//        <type>pending_trade</type>
//        <action>disallow</action>
//      </transaction>
//    </fantasy_content>
func (y *YahooConfig) AllowTrade() {
	// PUT

}
func (y *YahooConfig) DisallowTrade() {
	// PUT

}

// Voting Against Trades
// If there are accepted trades in your league waiting to be processed, which
// you can get by asking the transactions collection for all pending trades for
// your team, and you’re a manager in a league that allows managers to vote
// against trades, you can choose to vote against the trade. The input XML
// should look like:
//     <?xml version='1.0'?>
//     <fantasy_content>
//       <transaction>
//         <transaction_key>248.l.55438.pt.11</transaction_key>
//         <type>pending_trade</type>
//         <action>vote_against</action>
//         <voter_team_key>248.l.55438.t.2</voter_team_key>
//       </transaction>
//     </fantasy_content>
func (y *YahooConfig) VoteDownTrade() {
	// PUT
}

// DELETE
// Using DELETE, you may cancel any pending waiver claim or proposed trade. The
// URL for DELETEing a transaction resource is:
// http://fantasysports.yahooapis.com/fantasy/v2/transaction/
//
// You can only DELETE transactions of the types waiver or pending_trade if the
// pending trade has not yet been accepted.
func (y *YahooConfig) DeleteWaiver() {
	// DELETE
}
func (y *YahooConfig) DeletePendingTrade() {
	// DELETE
}

// Transactions collection
// With the Transactions API, you can obtain information via GET from a
// collection of transactions simultaneously. The transactions collection is
// qualified in the URI by a particular league. Each element beneath the
// Transactions Collection will be a Transaction Resource
//
// You can also POST to the API to perform operations like adding and/or
// dropping players to/from a team and proposing trades.
type TransactionCollection struct{}

// GetTransactionCollection
//
// URI:http://fantasysports.yahooapis.com/fantasy/v2/league//transactions
// Description:Fetch all completed transactions within a league.
// Sample:http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/transactions
//
// URI:http://fantasysports.yahooapis.com/fantasy/v2/transactions;transaction_keys=,{transaction_key2}
// Description:Fetch specific transactions {transaction_key1} and {transaction_key2}
// Sample:http://fantasysports.yahooapis.com/fantasy/v2/transactions;transaction_keys=223.l.431.tr.26,223.l.431.tr.27
//
// URI:http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys={league_key1},{league_key2}/transactions
// Description:Fetch all completed transactions of the leagues {league_key1} and {league_key2}
// Sample:http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=223.l.431,223.l.21821/transactions
//
// URI:http://fantasysports.yahooapis.com/fantasy/v2/league/{league_key}/transactions;types=waiver,pending_trade;team_key=
// Description:Fetch all pending trades and waivers relevant to the particular team.
// Sample:http://fantasysports.yahooapis.com/fantasy/v2/league/223.l.431/transactions;types=waiver,pending_trade;team_key=223.l.431.t.1
//
// Any sub-resource valid for a transaction is a valid sub-resource under the
// transactions collection.
// Any sub-resource for a collection of transactions is extracted using a URI
// like:
//   /transactions/{sub_resource}
//  OR
//   /transactions;transaction_keys={transaction_key1},{transaction_key2}/{sub_resource}
//
// Multiple sub-resources can be extracted from transactions in the same URI
// using a format like:
//   /transactions;out={sub_resource_1},{sub_resource_2}
//  OR
//   /transactions;transaction_keys={transaction_key1},{transaction_key2};out={sub_resource_1},{sub_resource_2}
//
// Filters
// The transactions collection can have filters such as the following to obtain
// a subset of a transactions collection that satisfy the filtering condition.
// These filters can be combined to obtain a more restricted list of
// transactions.
//
//
func (y *YahooConfig) GetTransactionCollection() *TransactionCollection {
	return nil
}

/*

Any sub-resource valid for a transaction is a valid sub-resource under the transactions collection.


Filters¶



Filter parameter	Filter parameter values	Usage
type	add,drop,commish,trade	/transactions;type=add
types	Any valid types	/transactions;types=add,trade
team_key	A team_key within the league	/transactions;team_key=257.l.193.t.1
type with team_key	waiver,pending_trade	You can only use these options when also providing the team_key, ie /transactions;team_key=257.l.193.t.1;type=waiver
count	Any integer greater than 0	/transactions;count=5
*/


// POST
// Using POST, players can be added and/or dropped from a team, or trades can be
// proposed. The URI for POSTing to transactions collection is:
// http://fantasysports.yahooapis.com/fantasy/v2/league//transactions
//
// Adding/Dropping Players
// The input XML format for a POST request to the transactions API for adding a
// player is:
//     <fantasy_content>
//       <transaction>
//         <type>add</type>
//         <player>
//           <player_key>{player_key}</player_key>
//           <transaction_data>
//             <type>add</type>
//             <destination_team_key>{team_key}</destination_team_key>
//           </transaction_data>
//         </player>
//       </transaction>
//     </fantasy_content>
//
// The input XML format for a POST request to the transactions API for dropping
// a player is:
//     <fantasy_content>
//       <transaction>
//         <type>drop</type>
//         <player>
//           <player_key>{player_key}</player_key>
//           <transaction_data>
//             <type>drop</type>
//             <source_team_key>{team_key}</source_team_key>
//           </transaction_data>
//         </player>
//       </transaction>
//     </fantasy_content>


// The input XML format for a POST request to the transactions API for replacing
// one player with another player in a team is:
//     <fantasy_content>
//       <transaction>
//         <type>add/drop</type>
//         <players>
//           <player>
//             <player_key>{player_key}</player_key>
//             <transaction_data>
//               <type>add</type>
//               <destination_team_key>{team_key}</destination_team_key>
//             </transaction_data>
//           </player>
//           <player>
//             <player_key>{player_key}</player_key>
//             <transaction_data>
//               <type>drop</type>
//               <source_team_key>{team_key}</source_team_key>
//             </transaction_data>
//           </player>
//         </players>
//       </transaction>
//     </fantasy_content>
//
// You may also add players that are currently on waivers – the players will not
// be immediately added to your team, but rather, you will be returned back a
// waiver claim that will be processed at some point in the future. Various
// league rules will control in which conditions you will actually receive the
// player, in the case that multiple teams have placed waiver claims.
//
// If you are placing a waiver claim in a league that uses FAAB, you may add
// that to the XML that you POST:
//     <?xml version='1.0'?>
//     <fantasy_content>
//       <transaction>
//         <type>add/drop</type>
//         <faab_bid>25</faab_bid>
//         <players>
//           <player>
//             <player_key>238.p.5484</player_key>
//             <transaction_data>
//               <type>add</type>
//               <destination_team_key>238.l.627060.t.6</destination_team_key>
//             </transaction_data>
//           </player>
//           <player>
//             <player_key>238.p.6327</player_key>
//             <transaction_data>
//               <type>drop</type>
//               <destination_team_key>238.l.627060.t.6</destination_team_key>
//             </transaction_data>
//           </player>
//         </players>
//       </transaction>
//     </fantasy_content>
//
// Once you have a waiver claim transaction, you may also edit the waiver
// priority or FAAB bid, or cancel the waiver entirely.
//
// Proposing Trades
// The input XML format for a POST request to the transactions API for proposing
// a trade is:
//     <?xml version='1.0'?>
//     <fantasy_content>
//       <transaction>
//         <type>pending_trade</type>
//         <trader_team_key>248.l.55438.t.11</trader_team_key>
//         <tradee_team_key>248.l.55438.t.4</tradee_team_key>
//         <trade_note>Yo yo yo yo yo!!!</trade_note>
//         <players>
//           <player>
//             <player_key>248.p.4130</player_key>
//             <transaction_data>
//               <type>pending_trade</type>
//               <source_team_key>248.l.55438.t.11</source_team_key>
//               <destination_team_key>248.l.55438.t.4</destination_team_key>
//             </transaction_data>
//           </player>
//           <player>
//             <player_key>248.p.2415</player_key>
//             <transaction_data>
//               <type>pending_trade</type>
//               <source_team_key>248.l.55438.t.4</source_team_key>
//               <destination_team_key>248.l.55438.t.11</destination_team_key>
//             </transaction_data>
//           </player>
//         </players>
//       </transaction>
//     </fantasy_content>

// Once you have a pending trade transaction, you may accept, reject, allow/
// disallow, or vote against the trade (depending on which role you have in the
// league). You may also cancel the trade.


// User resource
// With the User API, you can retrieve fantasy information for a particular
// Yahoo! user. Most usefully, you can see which games a user is playing, and
// which leagues they belong to and teams that they own within those games.
// Because you can currently only view user information for the logged in user,
// you would generally want to use the Users collection, passing along the
// use_login flag, instead of trying to request a User resource directly from
// the URI.
type UserResource struct {
}

// GetUserResource
// It is generally recommended that you instead use the Users collection,
// passing along the use_login flag.
//
// Sub-resources¶
// Name:
// Description: Fetch the Games in which the user has played. Additionally
//              accepts flags is_available to only return available games.
// URI:         /fantasy/v2/;use_login=1/games
// Sample:      http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games
//
// Name:        /
// Description: Fetch leagues that the user belongs to in one or more games. The leagues will be scoped to the user. This will throw an error if any of the specified games do not support league sub-resources.
// URI:         /fantasy/v2/;use_login=1/games;game_keys=,{game_key2}/leagues
// Sample:      http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys=223/leagues
//
// Name:
// Description: Fetch teams owned by the user in one or more games. The teams
//              will be scoped to the user. This will throw an error if any of
//              the specified games do not support team sub-resources.
// URI:         /fantasy/v2/;use_login=1/games;game_keys=,{game_key2}/teams
// Sample:      http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys=223/teams
func (y *YahooConfig) GetUserResource() *UserResource {
	// GET
	panic("Not Implemented")
	return nil
}

// Users collection
// With the Users API, you can obtain information from a collection of users
// simultaneously. Each element beneath the Users Collection will be a User
// Resource
type UserCollection struct {
	Body string
}

// Retrieve User Collection
// URIs
//
//   URI:
//   Description:Fetch user information of the logged-in user.
//   Sample:http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1
//
// Any sub-resource valid for a user is a valid sub-resource under the users collection.
// Any sub-resource for a collection of users is extracted using a URI like:
//     /users;use_login=1/{sub_resource}
// Multiple sub-resources can be extracted from users in the same URI using a format like:
//     /users;use_login=1;out={sub_resource_1},{sub_resource_2}
//     /users;field={field_name1},{field_name2}
func (y *YahooConfig) GetUserCollection(r *http.Request) *UserCollection {
	session, err := y.SessionStore.Get(r, "session-name")
	if err != nil {
		log.Println(err.Error(), 500)
		return nil
	}

  tok := session.Values["token"].(oauth2.Token)
  client := y.conf.Client(oauth2.NoContext, &tok)

	var userCollection UserCollection

	res, err := client.Get("http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "%s", body)
	userCollection.Body = string(body)

	// var animals []Animal
	// err := json.Unmarshal(jsonBlob, &animals)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// fmt.Printf("%+v", animals)

	return &userCollection
}
