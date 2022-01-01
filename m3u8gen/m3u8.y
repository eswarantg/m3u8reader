%{
package m3u8gen

import "github.com/eswarantg/m3u8reader"
import "time"

func TagName(token int) string {
      const TOKSTART = 4
      token -= TAG_FIRST + 1
      token += TOKSTART
      return yyToknames[token]
}
func AttrName(token int) string {
      const TOKSTART = 4
      token -= ATTR_FIRST + 1
      token += TOKSTART
      return yyToknames[token]
}

func setResult(l yyLexer, v *m3u8reader.M3U8) {
  l.(*Lexer).parseResult = v
}

%}

%union  {
    i int;
    i64 int64;
    f float64;
    s string;
    r string;
    t time.Time;
    val interface{}
    entry m3u8reader.M3U8Entry
    manifest m3u8reader.M3U8
}

%token TAG_FIRST
%token <i> TAG_EXTM3U
%token <i> TAG_EXT_X_VERSION
%token <i> TAG_EXT_X_INDEPENDENT_SEGMENTS
%token <i> TAG_EXT_X_MEDIA
%token <i> TAG_EXT_X_STREAM_INF
%token <i> TAG_EXT_X_TARGETDURATION
%token <i> TAG_EXT_X_SERVER_CONTROL
%token <i> TAG_EXT_X_PART_INF
%token <i> TAG_EXT_X_MEDIA_SEQUENCE
%token <i> TAG_EXT_X_SKIP
%token <i> TAG_EXTINF
%token <i> TAG_EXT_X_PROGRAM_DATE_TIME
%token <i> TAG_EXT_X_PART
%token <i> TAG_EXT_X_PRELOAD_HINT
%token <i> TAG_EXT_X_RENDITION_REPORT
%token <i> TAG_EXT_X_MAP

%token <i> COMMA
%token <i> EQUALTO
%token <s> SECONDLINEVALUE

%token ATTR_FIRST
%token <i> ATTR_BANDWIDTH
%token <i> ATTR_AVERAGE_BANDWIDTH
%token <i> ATTR_RESOLUTION
%token <i> ATTR_FRAME_RATE
%token <i> ATTR_CODECS
%token <i> ATTR_AUDIO
%token <i> ATTR_TYPE
%token <i> ATTR_GROUP_ID
%token <i> ATTR_NAME
%token <i> ATTR_DEFAULT
%token <i> ATTR_AUTOSELECT
%token <i> ATTR_LANGUAGE
%token <i> ATTR_CHANNELS
%token <i> ATTR_URI
%token <i> ATTR_CAN_BLOCK_RELOAD
%token <i> ATTR_CAN_SKIP_UNTIL
%token <i> ATTR_PART_HOLD_BACK
%token <i> ATTR_PART_TARGET
%token <i> ATTR_SKIPPED_SEGMENTS
%token <i> ATTR_DURATION
%token <i> ATTR_INDEPENDENT
%token <i> ATTR_LAST_MSN
%token <i> ATTR_LAST_PART

%token <s> ATTRKEY

%token <i64> INTEGERVAL
%token <f> FLOATVAL
%token <s> STRINGVAL
%token <r> RESOLUTIONVAL
%token <t> TIMEVAL

%type <i> ATTRTOKEN
%type <val> VALUE
%type <entry> ATTRLIST
%type <entry> ATTRANDVAL 
%type <entry> entry
%type <manifest> entries 
%type <manifest> manifest 

%start manifest

%%

manifest:  TAG_EXTM3U entries { setResult(yylex, &$$) }

entries :  entry  { $$.PostRecordEntry($1) }
        |  entries entry { $$.PostRecordEntry($2) }

entry : TAG_EXT_X_VERSION INTEGERVAL { $$.Tag = TagName($1);  $$.StoreKV("#",$2) } 
      | TAG_EXT_X_STREAM_INF ATTRLIST SECONDLINEVALUE { $$.Tag = TagName($1);  $$.StoreKV("#",$3) } 
      | TAG_EXT_X_INDEPENDENT_SEGMENTS { $$.Tag = TagName($1) } 
      | TAG_EXT_X_MEDIA ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXT_X_TARGETDURATION INTEGERVAL { $$.Tag = TagName($1);  $$.StoreKV("#",$2) } 
      | TAG_EXT_X_SERVER_CONTROL ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXT_X_PART_INF ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXT_X_MEDIA_SEQUENCE INTEGERVAL { $$.Tag = TagName($1);  $$.StoreKV("#",$2) } 
      | TAG_EXT_X_SKIP ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXTINF FLOATVAL COMMA SECONDLINEVALUE { $$.Tag = TagName($1);  $$.StoreKV("#",$2) ; $$.StoreKV("URI",$4) } 
      | TAG_EXTINF INTEGERVAL COMMA SECONDLINEVALUE { $$.Tag = TagName($1);  $$.StoreKV("#",float64($2)) ; $$.StoreKV("URI",$4) } 
      | TAG_EXT_X_PROGRAM_DATE_TIME TIMEVAL { $$.Tag = TagName($1);  $$.StoreKV("#",$2) } 
      | TAG_EXT_X_PART ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXT_X_PRELOAD_HINT ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXT_X_RENDITION_REPORT ATTRLIST { $$.Tag = TagName($1) } 
      | TAG_EXT_X_MAP ATTRLIST { $$.Tag = TagName($1) } 

ATTRLIST : ATTRANDVAL
         | ATTRLIST COMMA ATTRANDVAL 

ATTRANDVAL : ATTRTOKEN VALUE { $$.StoreKV(AttrName($1),$2) } 
           | ATTRKEY VALUE { $$.StoreKV($1,$2) } 

ATTRTOKEN : ATTR_BANDWIDTH { $$ = $1 }
          | ATTR_AVERAGE_BANDWIDTH { $$ = $1 }
          | ATTR_RESOLUTION { $$ = $1 }
          | ATTR_FRAME_RATE { $$ = $1 }
          | ATTR_CODECS { $$ = $1 }
          | ATTR_AUDIO { $$ = $1 }
          | ATTR_TYPE { $$ = $1 }
          | ATTR_GROUP_ID { $$ = $1 }
          | ATTR_NAME { $$ = $1 }
          | ATTR_DEFAULT { $$ = $1 }
          | ATTR_AUTOSELECT { $$ = $1 }
          | ATTR_LANGUAGE { $$ = $1 }
          | ATTR_CHANNELS { $$ = $1 }
          | ATTR_URI { $$ = $1 }
          | ATTR_CAN_BLOCK_RELOAD { $$ = $1 }
          | ATTR_CAN_SKIP_UNTIL { $$ = $1 }
          | ATTR_PART_HOLD_BACK { $$ = $1 }
          | ATTR_PART_TARGET { $$ = $1 }
          | ATTR_SKIPPED_SEGMENTS { $$ = $1 }
          | ATTR_DURATION { $$ = $1 }
          | ATTR_INDEPENDENT { $$ = $1 }
          | ATTR_LAST_MSN { $$ = $1 }
          | ATTR_LAST_PART { $$ = $1 }

VALUE : INTEGERVAL { $$ = $1 }
      | FLOATVAL { $$ = $1 }
      | STRINGVAL { $$ = $1 }
      | RESOLUTIONVAL { $$ = $1 }
      | TIMEVAL { $$ = $1 }
%%
