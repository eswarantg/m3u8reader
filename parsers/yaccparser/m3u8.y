%{
package yaccparser

import "time"
import "github.com/eswarantg/m3u8reader/common"
import "github.com/eswarantg/m3u8reader/parsers"

func tokenIdToTagId(token int) common.TagId {
	return common.TagId(token - tag_FIRST)
}
func attrTokenToTagId(token int) common.AttrId {
	return common.AttrId(token - token_ATTR_FIRST - 1)
}

func getHandler(l yyLexer) parsers.M3u8Handler {
	var lexer *Lexer
	var ok bool
	var obj parsers.M3u8Handler
	lexer, ok = l.(*Lexer)
	if !ok {
		panic("\nunknown lexer")
	}
	obj, ok = lexer.parseResult.(parsers.M3u8Handler)
	if !ok {
		panic("\nunknown object")
	}
	if obj == nil {
		panic("\nnil object")
	}
	return obj
}

%}

%union  {
    i int;
    i64 int64;
    f float64;
    s string;
    r string;
    t time.Time;
    kv keyValuePair;
    val interface{}
    kvpairs keyValuePairs
    entry accEntry
    hdlr parsers.M3u8Handler
}

/* 
Important to have this sequence same as the common.TagId sequence 
Also the m3u8.nex file needs to include the REGEX for the Tag
*/
%token tag_FIRST
%token <i> tag_EXTM3U
%token <i> tag_EXT_X_VERSION
%token <i> tag_EXT_X_INDEPENDENT_SEGMENTS
%token <i> tag_EXT_X_MEDIA
%token <i> tag_EXT_X_STREAM_INF
%token <i> tag_EXT_X_TARGETDURATION
%token <i> tag_EXT_X_SERVER_CONTROL
%token <i> tag_EXT_X_PART_INF
%token <i> tag_EXT_X_MEDIA_SEQUENCE
%token <i> tag_EXT_X_SKIP
%token <i> tag_EXTINF
%token <i> tag_EXT_X_PROGRAM_DATE_TIME
%token <i> tag_EXT_X_PART
%token <i> tag_EXT_X_PRELOAD_HINT
%token <i> tag_EXT_X_RENDITION_REPORT
%token <i> tag_EXT_X_MAP
/*
TBD
	"EXT-X-I-FRAME-STREAM-INF",
	"EXT-X-DISCONTINUITY",
	"EXT-X-ENDLIST",
	"EXT-X-PLAYLIST-TYPE",
	"EXT-X-BYTERANGE",
	"EXT-X-KEY",
	"EXT-X-DATERANGE",
	"EXT-X-DISCONTINUITY-SEQUENCE",
	"EXT-X-I-FRAMES-ONLY",
	"EXT-X-SESSION-KEY",
	"EXT-X-SESSION-DATA",
	"EXT-X-START",
*/

%token <i> token_COMMA
%token <s> token_SECONDLINEVALUE

/* 
Important to have this sequence same as the common.AttrId sequence 
Also the m3u8.nex file needs to include the REGEX for the Attribute
*/
%token token_ATTR_FIRST
%token <i> token_ATTR_BANDWIDTH
%token <i> token_ATTR_AVERAGE_BANDWIDTH
%token <i> token_ATTR_RESOLUTION
%token <i> token_ATTR_FRAME_RATE
%token <i> token_ATTR_CODECS
%token <i> token_ATTR_AUDIO
%token <i> token_ATTR_TYPE
%token <i> token_ATTR_GROUP_ID
%token <i> token_ATTR_NAME
%token <i> token_ATTR_DEFAULT
%token <i> token_ATTR_AUTOSELECT
%token <i> token_ATTR_LANGUAGE
%token <i> token_ATTR_CHANNELS
%token <i> token_ATTR_URI
%token <i> token_ATTR_CAN_BLOCK_RELOAD
%token <i> token_ATTR_CAN_SKIP_UNTIL
%token <i> token_ATTR_PART_HOLD_BACK
%token <i> token_ATTR_PART_TARGET
%token <i> token_ATTR_SKIPPED_SEGMENTS
%token <i> token_ATTR_DURATION
%token <i> token_ATTR_INDEPENDENT
%token <i> token_ATTR_LAST_MSN
%token <i> token_ATTR_LAST_PART

/* TBD
	"METHOD",
	"IV",
	"KEYFORMAT",
	"KEYFORMATVERSIONS",
	"BYTERANGE",
	"ID",
	"CLASS",
	"START-DATE",
	"END-DATE",
	"PLANNED-DURATION",
	"SCTE35-CMD",
	"SCTE35-OUT",
	"SCTE35-IN",
	"END-ON-NEXT",
	"ASSOC-LANGUAGE",
	"FORCED",
	"INSTREAM-ID",
	"CHARACTERISTICS",
	"HDCP-LEVEL",
	"VIDEO",
	"SUBTITLES",
	"CLOSED-CAPTIONS",
	"TIME-OFFSET",
	"PRECISE",
	"DATA-ID",
	"VALUE",
*/


%token <s> token_ATTRKEY

%token <i64> token_INTEGERVAL
%token <f> token_FLOATVAL
%token <s> token_STRINGVAL
%token <r> token_RESOLUTIONVAL
%token <t> token_TIMEVAL

%type <i> ATTRTOKEN
%type <val> VALUE
%type <kv> ATTRANDVAL 
%type <kvpairs> ATTRLIST
%type <entry> entry
%type <hdlr> entries

%start manifest

%%

manifest: entries

entries :  entry  { if $$ == nil { $$ = getHandler(yylex) }; $$.PostRecord($1.tag, $1.kvs); $1.clear("ENTRY1") }
        |  entries entry { if $$ == nil { $$ = getHandler(yylex) }; $$.PostRecord($2.tag, $2.kvs); $2.clear("ENTRY2") }

entry : tag_EXTM3U { $$.tag = tokenIdToTagId($1); } 
      | tag_EXT_X_VERSION token_INTEGERVAL { $$.tag = tokenIdToTagId($1); $$.storeKVDebug("EXT_X_VERSION",common.INTUnknownAttr, $2) }
      | tag_EXT_X_STREAM_INF ATTRLIST token_SECONDLINEVALUE { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_STREAM_INF_1", $2); $2.clear("EXT_X_STREAM_INF_1"); $$.storeKVDebug("EXT_X_STREAM_INF_2",common.INTUnknownAttr, $3); }
      | tag_EXT_X_INDEPENDENT_SEGMENTS { $$.tag = tokenIdToTagId($1) } 
      | tag_EXT_X_MEDIA ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_MEDIA",$2); $2.clear("EXT_X_MEDIA"); } 
      | tag_EXT_X_TARGETDURATION token_INTEGERVAL { $$.tag = tokenIdToTagId($1);  $$.storeKVDebug("EXT_X_TARGETDURATION",common.INTUnknownAttr,$2) } 
      | tag_EXT_X_SERVER_CONTROL ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_SERVER_CONTROL",$2); $2.clear("EXT_X_SERVER_CONTROL"); } 
      | tag_EXT_X_PART_INF ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_PART_INF",$2); $2.clear("EXT_X_PART_INF"); } 
      | tag_EXT_X_MEDIA_SEQUENCE token_INTEGERVAL { $$.tag = tokenIdToTagId($1);  $$.storeKVDebug("EXT_X_MEDIA_SEQUENCE",common.INTUnknownAttr,$2) } 
      | tag_EXT_X_SKIP ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_SKIP",$2); $2.clear("EXT_X_SKIP"); } 
      | tag_EXTINF token_FLOATVAL token_COMMA token_SECONDLINEVALUE { $$.tag = tokenIdToTagId($1);  $$.storeKVDebug("EXTINF_FL_1",common.INTUnknownAttr,$2) ; $$.storeKVDebug("EXTINF_FL_2",common.M3U8Uri,$4) } 
      | tag_EXTINF token_INTEGERVAL token_COMMA token_SECONDLINEVALUE { $$.tag = tokenIdToTagId($1);  $$.storeKVDebug("EXTINF_INT_1",common.INTUnknownAttr,float64($2)) ; $$.storeKVDebug("EXTINF_INT_2",common.M3U8Uri,$4) } 
      | tag_EXT_X_PROGRAM_DATE_TIME token_TIMEVAL { $$.tag = tokenIdToTagId($1); $$.storeKVDebug("EXT_X_PROGRAM_DATE_TIME",common.INTUnknownAttr,$2) } 
      | tag_EXT_X_PART ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_PART",$2); $2.clear("EXT_X_PART");   } 
      | tag_EXT_X_PRELOAD_HINT ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_PRELOAD_HINT",$2); $2.clear("EXT_X_PRELOAD_HINT");  } 
      | tag_EXT_X_RENDITION_REPORT ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_RENDITION_REPORT",$2); $2.clear("EXT_X_RENDITION_REPORT");  } 
      | tag_EXT_X_MAP ATTRLIST { $$.tag = tokenIdToTagId($1); $$.assignKVPS("EXT_X_MAP",$2); $2.clear("EXT_X_MAP"); } 

ATTRLIST : ATTRANDVAL { $$.storeKVDebug("ATTRANDVAL_1",$1.k, $1.v) }
         | ATTRLIST token_COMMA ATTRANDVAL { $1.storeKVDebug("ATTRANDVAL_2",$3.k, $3.v); $$ = $1 } 

ATTRANDVAL : ATTRTOKEN VALUE { $$.k = attrTokenToTagId($1); $$.v=$2 } 
           | token_ATTRKEY VALUE { $$.k = common.AttrToAttrId[$1]; $$.v=$2 } 

ATTRTOKEN : token_ATTR_BANDWIDTH { $$ = $1 }
          | token_ATTR_AVERAGE_BANDWIDTH { $$ = $1 }
          | token_ATTR_RESOLUTION { $$ = $1 }
          | token_ATTR_FRAME_RATE { $$ = $1 }
          | token_ATTR_CODECS { $$ = $1 }
          | token_ATTR_AUDIO { $$ = $1 }
          | token_ATTR_TYPE { $$ = $1 }
          | token_ATTR_GROUP_ID { $$ = $1 }
          | token_ATTR_NAME { $$ = $1 }
          | token_ATTR_DEFAULT { $$ = $1 }
          | token_ATTR_AUTOSELECT { $$ = $1 }
          | token_ATTR_LANGUAGE { $$ = $1 }
          | token_ATTR_CHANNELS { $$ = $1 }
          | token_ATTR_URI { $$ = $1 }
          | token_ATTR_CAN_BLOCK_RELOAD { $$ = $1 }
          | token_ATTR_CAN_SKIP_UNTIL { $$ = $1 }
          | token_ATTR_PART_HOLD_BACK { $$ = $1 }
          | token_ATTR_PART_TARGET { $$ = $1 }
          | token_ATTR_SKIPPED_SEGMENTS { $$ = $1 }
          | token_ATTR_DURATION { $$ = $1 }
          | token_ATTR_INDEPENDENT { $$ = $1 }
          | token_ATTR_LAST_MSN { $$ = $1 }
          | token_ATTR_LAST_PART { $$ = $1 }

VALUE : token_INTEGERVAL { $$ = $1 }
      | token_FLOATVAL { $$ = $1 }
      | token_STRINGVAL { $$ = $1 }
      | token_RESOLUTIONVAL { $$ = $1 }
      | token_TIMEVAL { $$ = $1 }
%%
