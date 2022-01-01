%{
package m3u8reader

import "time"

func tokenIdToTagId(token int) TagId {
	return TagId(token - TAG_FIRST - 1)
}
func attrTokenToTagId(token int) AttrId {
	return AttrId(token - ATTR_FIRST - 1)
}

func getM3U8Store(l yyLexer) *M3U8 {
  var lexer *Lexer
  var ok bool
  var obj *M3U8
  lexer, ok = l.(*Lexer)
  if !ok {
        panic("unknown lexer")
  }
  obj, ok = lexer.parseResult.(*M3U8)
  if !ok {
        panic("unknown object")
  }
  if obj == nil {
        panic("nil object")
  }
  return obj
}

type KeyValuePair struct{
   k AttrId
   v interface{}
}





%}

%union  {
    i int;
    i64 int64;
    f float64;
    s string;
    r string;
    t time.Time;
    kv KeyValuePair;
    val interface{}
    entry M3U8Entry
    manifest *M3U8
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
%type <kv> ATTRANDVAL 
%type <entry> entry
%type <manifest> entries

%start manifest

%%

manifest: entries

entries :  entry  { if $$ == nil { $$=getM3U8Store(yylex) }; $$.postRecordEntry($1); $1=M3U8Entry{} }
        |  entries entry { if $$ == nil { $$=getM3U8Store(yylex) }; $1.postRecordEntry($2); $2=M3U8Entry{}; }

entry : TAG_EXTM3U { $$.Tag = tokenIdToTagId($1); } 
      | TAG_EXT_X_VERSION INTEGERVAL { $$.Tag = tokenIdToTagId($1);  $$.storeKV(INTUnknownAttr,$2) } 
      | TAG_EXT_X_STREAM_INF ATTRLIST SECONDLINEVALUE { $$.Tag = tokenIdToTagId($1); $2.storeKV(INTUnknownAttr,$3); $$.Values = $2.Values }
      | TAG_EXT_X_INDEPENDENT_SEGMENTS { $$.Tag = tokenIdToTagId($1) } 
      | TAG_EXT_X_MEDIA ATTRLIST { $$.Tag = tokenIdToTagId($1); $$.Values = $2.Values  } 
      | TAG_EXT_X_TARGETDURATION INTEGERVAL { $$.Tag = tokenIdToTagId($1);  $$.storeKV(INTUnknownAttr,$2) } 
      | TAG_EXT_X_SERVER_CONTROL ATTRLIST { $$.Tag = tokenIdToTagId($1) ; $$.Values = $2.Values } 
      | TAG_EXT_X_PART_INF ATTRLIST { $$.Tag = tokenIdToTagId($1) ; $$.Values = $2.Values } 
      | TAG_EXT_X_MEDIA_SEQUENCE INTEGERVAL { $$.Tag = tokenIdToTagId($1);  $$.storeKV(INTUnknownAttr,$2) } 
      | TAG_EXT_X_SKIP ATTRLIST { $$.Tag = tokenIdToTagId($1) ; $$.Values = $2.Values } 
      | TAG_EXTINF FLOATVAL COMMA SECONDLINEVALUE { $$.Tag = tokenIdToTagId($1);  $$.storeKV(INTUnknownAttr,$2) ; $$.storeKV(M3U8Uri,$4) } 
      | TAG_EXTINF INTEGERVAL COMMA SECONDLINEVALUE { $$.Tag = tokenIdToTagId($1);  $$.storeKV(INTUnknownAttr,float64($2)) ; $$.storeKV(M3U8Uri,$4) } 
      | TAG_EXT_X_PROGRAM_DATE_TIME TIMEVAL { $$.Tag = tokenIdToTagId($1);  $$.storeKV(INTUnknownAttr,$2) } 
      | TAG_EXT_X_PART ATTRLIST { $$.Tag = tokenIdToTagId($1); $$.Values = $2.Values  } 
      | TAG_EXT_X_PRELOAD_HINT ATTRLIST { $$.Tag = tokenIdToTagId($1) ; $$.Values = $2.Values } 
      | TAG_EXT_X_RENDITION_REPORT ATTRLIST { $$.Tag = tokenIdToTagId($1) ; $$.Values = $2.Values } 
      | TAG_EXT_X_MAP ATTRLIST { $$.Tag = tokenIdToTagId($1); $$.Values = $2.Values  } 

ATTRLIST : ATTRANDVAL { $$.storeKV($1.k, $1.v) }
         | ATTRLIST COMMA ATTRANDVAL { $1.storeKV($3.k, $3.v); $$ = $1 } 

ATTRANDVAL : ATTRTOKEN VALUE { $$.k = attrTokenToTagId($1); $$.v=$2 } 
           | ATTRKEY VALUE { $$.k = attrToAttrId[$1]; $$.v=$2 } 

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
