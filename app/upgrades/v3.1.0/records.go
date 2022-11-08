package v310

// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.42/365)*13+1) | floor) | tostring})' DAN.JSON > to_mint.json

// Slash was 5%
// Lost APR is 42% for 13 days

var recordsJSONString = `[
  {
    "address": "chihuahua1qf6c6crk6a2uyndt3tnd4xazucpzvyslrxm34s",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1q2t96wltse9s8ygl2a4e92tccmahrve6yut5xq",
    "amount": "9829876986"
  },
  {
    "address": "chihuahua1qtfp0glcmz3jrfr02q66d702da6lj7eeqlw2jy",
    "amount": "15462645"
  },
  {
    "address": "chihuahua1qvmpf9y7g8hqu68f5mq5w9fkdkjkv5jw6mzvru",
    "amount": "7894821043"
  },
  {
    "address": "chihuahua1q39p26ve7gv7r93hdxhuqvrg7c8e53krq064j3",
    "amount": "179721564"
  },
  {
    "address": "chihuahua1qjh6njg837sjyzgxfsjhzdzvcs90cudgchafzk",
    "amount": "50697197"
  },
  {
    "address": "chihuahua1qly96sy3yklkwv6cf4he5ves4qu8yn58ea7wn5",
    "amount": "1023122076"
  },
  {
    "address": "chihuahua1pzs425ltmdtrshh5ehc24l6xtwudp7l009tq34",
    "amount": "3760211120"
  },
  {
    "address": "chihuahua1pys48fndsf8552a35e0lqhv7487xnqjdf4qspg",
    "amount": "4463383"
  },
  {
    "address": "chihuahua1py7vnejer4auv8e0dyexvgkas4q08wppjyl99y",
    "amount": "101394394"
  },
  {
    "address": "chihuahua1pxnzj2nfffwtkts6h7hy40wzc4fflj6lndm6wz",
    "amount": "62809642999"
  },
  {
    "address": "chihuahua1pgj7lah4elg5clyxyxj7mu7cgtqnt06md2z73g",
    "amount": "1750043894"
  },
  {
    "address": "chihuahua1pwjc9nses6qzt73kvvm2uelzvfmz7gleffzpe2",
    "amount": "101394394520"
  },
  {
    "address": "chihuahua1pj9e89gc8dkvv72wh5qdsxwrnwzev8tvyzggsp",
    "amount": "17744019041"
  },
  {
    "address": "chihuahua1p45vamfvac0txmtd0hgq4ql3a6spf4s65anz7y",
    "amount": "512041"
  },
  {
    "address": "chihuahua1p6906xu66zet86hne5z3m6lrekc29gw8d6gvy8",
    "amount": "2259928962"
  },
  {
    "address": "chihuahua1p7lgulzc3z20yhjfa2tpe2rylwfxhaktnpg0u2",
    "amount": "13181271"
  },
  {
    "address": "chihuahua1zqjkr6qgu7zmvj6gwy7gmzedtqsrsrunaasep9",
    "amount": "0"
  },
  {
    "address": "chihuahua1zz2jf3kjcrm909r7g2qfjmcqyu4kj9asqe9pg0",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1z8kqeftwhfqy28md3xnvf9dr9the3etlrv2sd7",
    "amount": "6185058"
  },
  {
    "address": "chihuahua1z86328atna82gh36tppwxepu000jdyxdll593q",
    "amount": "12155160022"
  },
  {
    "address": "chihuahua1zg9swa5rslrxfykr7ltxyky6h03w5ax3cau39r",
    "amount": "55310642"
  },
  {
    "address": "chihuahua1zvk2aat2lkpxxdjvm0kknwf7hm4ymne8cqvzjs",
    "amount": "20278878904"
  },
  {
    "address": "chihuahua1zjfvjpsr47xepxtnaakdglh7v5zy8effx2ymah",
    "amount": "26365761"
  },
  {
    "address": "chihuahua1zhm4k3qs0tm20v656t7ckn5hysqjt788lcwmnn",
    "amount": "1013943"
  },
  {
    "address": "chihuahua1rq98esw4lgcekya6yyksktrppcwe2hhc2x9ysu",
    "amount": "5932256491"
  },
  {
    "address": "chihuahua1rqjrpkntamtyjf0yvswctmzdkqftf5z9xv0x52",
    "amount": "49196580"
  },
  {
    "address": "chihuahua1rywdk96c75hn4vx8aj3wsasdtlj0968chzt864",
    "amount": "14245912"
  },
  {
    "address": "chihuahua1ryafk0pmg92tpcvf75g393kc8d0ym89g6y4vty",
    "amount": "304183183"
  },
  {
    "address": "chihuahua1r92ggd6d6gmmnvewyx4rxx32unmtuywnqcx4we",
    "amount": "538116291"
  },
  {
    "address": "chihuahua1r8he2t2eurm60d2ad45x5xncfv7tt2rlmrq3jd",
    "amount": "71047123287"
  },
  {
    "address": "chihuahua1rg9v6vp4lq2s3uar2q4jzf04t4zuwvyr45yvvh",
    "amount": "50697197260"
  },
  {
    "address": "chihuahua1rnxcgeehy7mfpt2zfyrsqpfmqfq4nluuvcglm5",
    "amount": "6912723294"
  },
  {
    "address": "chihuahua1rk6lejeurkdjn6cvfx42gpwuhzuf3t20g7h7uj",
    "amount": "25373972602"
  },
  {
    "address": "chihuahua1rel0d5hlj6v9dune23svdk24maz9rl0wdlrdjw",
    "amount": "506971972"
  },
  {
    "address": "chihuahua1raacj600qfth8pf6ngdamfawy0rjdmvjlh2wv5",
    "amount": "101394394"
  },
  {
    "address": "chihuahua1ypk8p72lqfezjpypxwu6n69de6qzm2syukld38",
    "amount": "12672405"
  },
  {
    "address": "chihuahua1y9x0q6fvp8uenpjas2g2jqzk0f22ptxrs9dvxs",
    "amount": "101140908"
  },
  {
    "address": "chihuahua1ygfmz8wrkspye8hd9pyhz5dwexwpjx2wq2y5hq",
    "amount": "506"
  },
  {
    "address": "chihuahua1yfrj43etaqp26mqluuvlz3350hkyavdh7hwvs7",
    "amount": "2332071073"
  },
  {
    "address": "chihuahua1yn2rkp8k3txj3wh5ks234zql5ral8fk40500gf",
    "amount": "69427857850"
  },
  {
    "address": "chihuahua1yu3u7xguyn79xp4h08wc9za6k5zes3dzm4p8dw",
    "amount": "39687289759"
  },
  {
    "address": "chihuahua19rfqwv3equna5eyypmhrw7yvc0qz0y9mjjt8fd",
    "amount": "2581861319"
  },
  {
    "address": "chihuahua198uvyllue9hpawr8jpsmhwnkpq5765yvxe9u4f",
    "amount": "506971"
  },
  {
    "address": "chihuahua192vksf2x5cqz8csx2hw57ksqafvfxa6rsvskmc",
    "amount": "59989475330"
  },
  {
    "address": "chihuahua19tutwdztw6rty9dnz6awxw0dtnuztytj6wunam",
    "amount": "6621662328"
  },
  {
    "address": "chihuahua195z08s2qj0em8pzyyaqm5l29p0w9vrde83qx6m",
    "amount": "152243835616"
  },
  {
    "address": "chihuahua197l3uzss7eeral5va35kaj3p2fcz2gz5ayjm7u",
    "amount": "4157170"
  },
  {
    "address": "chihuahua1xp6j5kxgv5zl3qw5cw72s48cq3rdh6xpxzdvx0",
    "amount": "388796805"
  },
  {
    "address": "chihuahua1x8chd4ejed8zuxrnrwwf97yt2n0x44cgh624rj",
    "amount": "2534859863"
  },
  {
    "address": "chihuahua1xg3jhs6d0s56gkrzwqjrc26tuafxxu6wrwvhy3",
    "amount": "50752969"
  },
  {
    "address": "chihuahua1x2p00vfqlfwntrcfz3d9n74pyzdy808v3x77w3",
    "amount": "2156787671"
  },
  {
    "address": "chihuahua1xw79afn6ue65qa0f592snhu0usweugyx426mmd",
    "amount": "2218891184"
  },
  {
    "address": "chihuahua1xshcw4wjm3kkekgevg8yfewvk6dhnzsktz639u",
    "amount": "0"
  },
  {
    "address": "chihuahua1xs79nssadxszdaefjwsvd8wsea965sthax3znc",
    "amount": "2994025605"
  },
  {
    "address": "chihuahua1xjecqq4qakmv4lun6hgus27y9c5akvwkyzzpnh",
    "amount": "50697"
  },
  {
    "address": "chihuahua1x4856r4ha4u0zj76hva08lyql5hrx5pxa3x2vm",
    "amount": "20025392"
  },
  {
    "address": "chihuahua1xcu6k2pw5mye8nxkvu2m24lrgc0ae8vyfmlndt",
    "amount": "506971"
  },
  {
    "address": "chihuahua18zdpn5zm2gxup73dhmrmuc58rnjrxnvc9lragq",
    "amount": "8314340"
  },
  {
    "address": "chihuahua18rfdktgpxeepy52n57zmu64y6wz0pxqyavfm52",
    "amount": "8238294554"
  },
  {
    "address": "chihuahua189atcn4hr73q4w95ym9fnnt5e4y0fj6avghks2",
    "amount": "44721214436"
  },
  {
    "address": "chihuahua186qqhn9um6rksed7fqzg5uxhlfs099dkac62q6",
    "amount": "10033330218"
  },
  {
    "address": "chihuahua1gqcsp5u2j3uew074en3jv8ekcn9y0xss9dg5ey",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1gpkg38400ufhswzxcx99c2g7ft7n7fpgldck9u",
    "amount": "7312804599"
  },
  {
    "address": "chihuahua1gzxk0g6mt63k23fm694lmh4zjd6wnxp2xd3w5v",
    "amount": "12155160022"
  },
  {
    "address": "chihuahua1gynxckgz6agq05duvwcvn3enzw97xu6jyyqxyh",
    "amount": "4562747"
  },
  {
    "address": "chihuahua1gylzr0qr6ekq6cpwcjlttqtyyveen74grq7zzc",
    "amount": "463225223"
  },
  {
    "address": "chihuahua1gxr53f9jx5cxn3dmkzr2eq93g05yey3054psw0",
    "amount": "50747945205"
  },
  {
    "address": "chihuahua1g80l3p523v0qak88k39u8rxc3fq0jadrwzp39u",
    "amount": "2585557"
  },
  {
    "address": "chihuahua1g8lfd4q6xcuu24fuawwsaej72e4t2rpy2nk8js",
    "amount": "11580423539"
  },
  {
    "address": "chihuahua1gfxwv0cm9a6h63g48tufvz5grge3h57ymglru0",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1gvg43c4q0d8nyrpt0e0gdwqxpf4lh886rgj3ht",
    "amount": "16871683675"
  },
  {
    "address": "chihuahua1g0dgu4zt5p0l9qevagf7sgectqel06lwcsetz8",
    "amount": "0"
  },
  {
    "address": "chihuahua1g0hpuuquyk3t2hgtm74336jvpy6ssue7vgj5sq",
    "amount": "33293457807"
  },
  {
    "address": "chihuahua1gs4u2hrdqlzlz8f8nzpp55zd2l2jpz3vr7yj4d",
    "amount": "10754938569"
  },
  {
    "address": "chihuahua1g3s6z9gndyfd24rpvt9sefqm4syv2cn5eplkqs",
    "amount": "18732869352"
  },
  {
    "address": "chihuahua1ghfaapsausm6kn9ush6wa4ggw6thvjgwzvkvq8",
    "amount": "10282837282"
  },
  {
    "address": "chihuahua1g6mzmawu0ansnml9gj0ffcecvryjcsewtxyfnw",
    "amount": "95923155"
  },
  {
    "address": "chihuahua1fg38wdrxkwaslc4aw6j4egqum4u45guucx07zs",
    "amount": "15158461"
  },
  {
    "address": "chihuahua1fdfgyjjdx2xffzctjzgmvpp5spakj8gq2fa8an",
    "amount": "1051628869"
  },
  {
    "address": "chihuahua1f3latzuk0qhmpkmp37v4uhurmpu7wnv89n4efa",
    "amount": "51639942595"
  },
  {
    "address": "chihuahua1fc2w5avsmcd4xk4rxquhv3rzdxz6kkhwyv8nsv",
    "amount": "101862128181"
  },
  {
    "address": "chihuahua12qpgx8l2f7x7pyrtvs8ky2dacefj007cxrx6gk",
    "amount": "650698526"
  },
  {
    "address": "chihuahua12p5wvy4nyywz5wm5y29j3rggtupzet5csnlhyh",
    "amount": "12674299"
  },
  {
    "address": "chihuahua12p6m6eppmffardltlhesltcvppzqvy5hkhfvav",
    "amount": "330038"
  },
  {
    "address": "chihuahua12rpjd07fws4hf3evejhxxdhzm5rn9862nwldhd",
    "amount": "50697"
  },
  {
    "address": "chihuahua12vghx05jf5e2sxgfkdhktl3ch0e2k9ae6jc7v9",
    "amount": "509767602"
  },
  {
    "address": "chihuahua120p506mdu8epen7age0antzmqxfl76yyk7cgyj",
    "amount": "0"
  },
  {
    "address": "chihuahua120aa8xgcxw97t7stzj6xvwwfdqkru8ma6c3t29",
    "amount": "1890678282"
  },
  {
    "address": "chihuahua123scuwjkz9v778sruaymudvvsfgv0a8wwv4lj4",
    "amount": "49371299706"
  },
  {
    "address": "chihuahua1234eths8eaw2pqzcs3sms8rzewmhjkrh9cwj7m",
    "amount": "5171114"
  },
  {
    "address": "chihuahua12n44yj8q89zxne05rajxqceufjhw3wrep07pk0",
    "amount": "104473874886"
  },
  {
    "address": "chihuahua125tl3umwps4ulsnkuwffqnw8mgadnljl96zp5v",
    "amount": "1466715"
  },
  {
    "address": "chihuahua1tp7dgz8eh8dfk2awjjnjrpn60tu2759nj7y4al",
    "amount": "64486834"
  },
  {
    "address": "chihuahua1tz4ssc67mz09pyc5xtvt3zvkm8g6uggkkhyqqp",
    "amount": "1268"
  },
  {
    "address": "chihuahua1tym288c48szfqcerrp57cvg3xl9ka5uuqrgww8",
    "amount": "1420942465"
  },
  {
    "address": "chihuahua1t9nd2wcgjxkw7mvgzm7z70vl6jf408u7r9q9jy",
    "amount": "15224383561"
  },
  {
    "address": "chihuahua1t89er2ysz7r8ntgaqzs0np0u9gasdupltp9xnz",
    "amount": "111533"
  },
  {
    "address": "chihuahua1t8hgejlx54gcz43mcqfuz874yamgtmgpt6237c",
    "amount": "253485986"
  },
  {
    "address": "chihuahua1tvyem2ed7e7r8c93q5rse3ut5nwwssvvf87dae",
    "amount": "80564533715"
  },
  {
    "address": "chihuahua1tady5gfzzwlzcjkmn8suzxp0uxf54vvxrycy0z",
    "amount": "33483217416"
  },
  {
    "address": "chihuahua1tacx5lykjmvdd5afkgygmuyju8h7flp0gfufu9",
    "amount": "6232509645"
  },
  {
    "address": "chihuahua1tl9fwdu2r6vflvz4pae7nnd373qn5ymhveqjnq",
    "amount": "5490927671"
  },
  {
    "address": "chihuahua1v8hayns4zauhwvw07tjf67lknq4ptdq78avltp",
    "amount": "2801273634"
  },
  {
    "address": "chihuahua1vwr50lmgcz94asaw2k08wfz72rsccn23t3w8q0",
    "amount": "11051989"
  },
  {
    "address": "chihuahua1vcgxjeeyprn8ppwm4gy4zwgdd30078r29ple89",
    "amount": "12243031699"
  },
  {
    "address": "chihuahua1v6689fk7tckkpfunz3nx90680p6hsqqq40z2v4",
    "amount": "0"
  },
  {
    "address": "chihuahua1vasunpy8eyx29fd4pk9kvkl8je2hpfvyz5umsy",
    "amount": "0"
  },
  {
    "address": "chihuahua1v7svumfrrdzcpx64cynl0dzhtqrkgrf6xgnzxu",
    "amount": "4562747"
  },
  {
    "address": "chihuahua1dyfc8v0m0qmnr2pyyual26yfmefrkk2v6lfv03",
    "amount": "2889740"
  },
  {
    "address": "chihuahua1d97peegdkq7t4e9c2rsdz4yf67dl24zyvyz4kr",
    "amount": "18250991"
  },
  {
    "address": "chihuahua1dxg87yptew6hg6rr3lhfh4ltkhdhmdx8yqeuwm",
    "amount": "30367621"
  },
  {
    "address": "chihuahua1d0fx9xxktc2sguth7q02lgfa47tuztz94ykg59",
    "amount": "31383036794"
  },
  {
    "address": "chihuahua1ds0spl74p7hvxt039s9mazygrpg6f9ashh08xv",
    "amount": "253485986"
  },
  {
    "address": "chihuahua1d496c6hyj95vl0cluhfkwcl476fg7hp9celk67",
    "amount": "53280799310"
  },
  {
    "address": "chihuahua1dep3djza75nyz26m5zj73zfrgcl07lvqt94z2g",
    "amount": "15944861145"
  },
  {
    "address": "chihuahua1daw2sahp0qcqcccv75p39wccj3guv6hkf0qju2",
    "amount": "0"
  },
  {
    "address": "chihuahua1wrllfdudjx3v67nrc8gqrwsj0rq4c0ejy9u8rf",
    "amount": "164933191"
  },
  {
    "address": "chihuahua1wxfn9dwuvglgj0evxnytgzqxp53t3xad2u7zsz",
    "amount": "2534859"
  },
  {
    "address": "chihuahua1wd4sx5z36letuyj0zd50v0kjspvzf8t6p2988x",
    "amount": "70976076"
  },
  {
    "address": "chihuahua1w0cv3dd8ugd9q028kmr064x53ah6fjxmrjhuws",
    "amount": "1682589279"
  },
  {
    "address": "chihuahua1ws9jelyyetmtt2y7z8gtrfgfy3mtveuq3jr5x7",
    "amount": "74025508"
  },
  {
    "address": "chihuahua1wk7hnadvkuftfxdm8hjz38vtsq356tq3hpqvps",
    "amount": "253485986"
  },
  {
    "address": "chihuahua1wh3mzs6yv3qjq32jwsearwakp3pxsfsv2lxpug",
    "amount": "210900340"
  },
  {
    "address": "chihuahua10gxum6e3klx2swgsxnuedmfs6tfu5txr4y042h",
    "amount": "28977942739"
  },
  {
    "address": "chihuahua10gk4v3lt3643jd7e2msdvjudpfurdeg7vqzx4p",
    "amount": "278023429"
  },
  {
    "address": "chihuahua10fdwe698fl5jfuq38lh6tgs0f0jp9uak53pkf4",
    "amount": "0"
  },
  {
    "address": "chihuahua10wxn2lv29yqnw2uf4jf439kwy5ef00qdn93mry",
    "amount": "151542034162"
  },
  {
    "address": "chihuahua10en6xp30u2qk8flt4nl58etttgs9x0qnp3qyu5",
    "amount": "20268739464"
  },
  {
    "address": "chihuahua10m0k4a3tudk6mglc9nez5w9yjvgkk9kuvl3j54",
    "amount": "0"
  },
  {
    "address": "chihuahua107m09gnrav7sfk8awf0fdnd5dn4fdlaulxhq5f",
    "amount": "3407218259"
  },
  {
    "address": "chihuahua1sx7kdx4zxdsss9reggd3vvdqp2462ckpy64ws4",
    "amount": "12155160022"
  },
  {
    "address": "chihuahua1swj9w30e8lzc5uxewe83nkun9305r062vhfcmv",
    "amount": "50697"
  },
  {
    "address": "chihuahua1s3y5g0nctfv6zvlwnw45gtk222x3ya2p3mnndk",
    "amount": "167300750"
  },
  {
    "address": "chihuahua1sjeudgjzaxe8nvghv3adua0w3mg890d5vnywzl",
    "amount": "253485986"
  },
  {
    "address": "chihuahua1sntx5sf6w933f8pytpr2dh9sxns0u23mfzz7v4",
    "amount": "29687547"
  },
  {
    "address": "chihuahua1s49ytjsawm4usq90r6rmazy8l2dx3l2v7x7tl4",
    "amount": "1589965500"
  },
  {
    "address": "chihuahua1s468q4uf98jzrtruvprcw6pqgdphakq3kxjcax",
    "amount": "116603553"
  },
  {
    "address": "chihuahua1smqhaen3y9l67587hwjmnezyngkyy6kzx7n6p7",
    "amount": "2719701844"
  },
  {
    "address": "chihuahua1saz2nvwyr87aqu60qc4gmnfl3s4s0qquknxt8s",
    "amount": "1016185"
  },
  {
    "address": "chihuahua1s7jhu47u2ukrcy0tpayrksmakkdkgjmv8z3rxu",
    "amount": "12053867022"
  },
  {
    "address": "chihuahua1s7k8uwcvzym8ltul405p4gyvgn8gq6q5rfwrs9",
    "amount": "1267429931"
  },
  {
    "address": "chihuahua13q6wx8vssvg2levkhwcynz29dms8m4q22wyhtw",
    "amount": "47959548"
  },
  {
    "address": "chihuahua13yc4kk6nytl0fzxu74md5s2eg3kguhwt3jhmpc",
    "amount": "12927785301"
  },
  {
    "address": "chihuahua132ufxdddhscavk2zglnd7465wqqlc6rykplrm3",
    "amount": "34762386713"
  },
  {
    "address": "chihuahua13vcc39q4v8aupx0ccxxhscg0l9karkjwgvc03c",
    "amount": "1394274420552"
  },
  {
    "address": "chihuahua133yfg2ufej32g6svv4tenqlp8m843s724v0psj",
    "amount": "14499398"
  },
  {
    "address": "chihuahua13306hrp8e9y9z8lv5h2eccupffj8kw95uqc5j6",
    "amount": "50752969"
  },
  {
    "address": "chihuahua13ketvyvzqxh6rq08qez3dvekd02uxll60ca2zk",
    "amount": "61310672710"
  },
  {
    "address": "chihuahua13h05p4kxh5mddd82g7mgngzu247lsyqkh0s345",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1jqz5panepz7gvdm3j4tq0h4rf5zdtq04g7sfsk",
    "amount": "20530643638"
  },
  {
    "address": "chihuahua1jyrtdm6xkaq2y6h5nkzrr2cxmn22fsq6ajwe5e",
    "amount": "7351093"
  },
  {
    "address": "chihuahua1j9cuhd3lrwsm98dj2hs6wf8t90580c56m3qlr4",
    "amount": "994706163"
  },
  {
    "address": "chihuahua1j9en47cu0k29j8xgyfxltruneeaw2rznjqz9u0",
    "amount": "101394394"
  },
  {
    "address": "chihuahua1jgfc0zgcdah5hnn6rep5dnu6u90alwalpaddpe",
    "amount": "5000843157"
  },
  {
    "address": "chihuahua1j2hzlc906peumvhc7tnmwhd587jtzsqqshky2w",
    "amount": "15209159"
  },
  {
    "address": "chihuahua1jdgyz058su4mhp3v0ejf8jgegdwage2u8208us",
    "amount": "16899048854"
  },
  {
    "address": "chihuahua1jstn2l0up0wlm5wx64dve6w6zxc5lglvh9t2wg",
    "amount": "6854427100"
  },
  {
    "address": "chihuahua1jjftv5f3emy4p4rfep0593q5jy4x6ddurk6uku",
    "amount": "8722452788"
  },
  {
    "address": "chihuahua1j5m2qugygwgwszqzjzcaz3e02y9uygtnu2wccz",
    "amount": "76045795"
  },
  {
    "address": "chihuahua1jeup594t8m8m2zg3j42ss3jm0c07cwh0k75xhl",
    "amount": "15856268926"
  },
  {
    "address": "chihuahua1judj83pu32eufkav82xuwhu04kkcqggns4f7a3",
    "amount": "35488038"
  },
  {
    "address": "chihuahua1ju3z4nef3ppxz8qed0cqcyytw0s6js5z7j8rw7",
    "amount": "7191803008"
  },
  {
    "address": "chihuahua1ja2w03e9wvd9vf9kc74n09q3yf76w5yf6wwqws",
    "amount": "25348598630"
  },
  {
    "address": "chihuahua1nzf6827gpnpfdzcfvultvr2f772zfzvgqgsvsk",
    "amount": "6235755"
  },
  {
    "address": "chihuahua1nyw2nhwttjld2natpxs6w4h76sg3r94mgzr0pq",
    "amount": "0"
  },
  {
    "address": "chihuahua1nyw5g9yjtcrcsvk4j5tzrk5fh9yjhshl7x0jau",
    "amount": "496832533"
  },
  {
    "address": "chihuahua1nyumr38hm3cwun74szrazdua2lu5u2wp5gj60f",
    "amount": "21546308"
  },
  {
    "address": "chihuahua1nt66tcttrm2htp53ppxs3e2787t95uplxtn0dc",
    "amount": "200631882"
  },
  {
    "address": "chihuahua1nwkvkm2rwfm9wve2cl3h8gvs7cvcfhj2sk728v",
    "amount": "10139439452"
  },
  {
    "address": "chihuahua1n0w4stkdkxfd5pcm46h4duh4clewnfsc32gs4t",
    "amount": "5712500048"
  },
  {
    "address": "chihuahua1nckqhfp8k67qzvats2slqvtaf3kynz66pvhjqh",
    "amount": "771966222"
  },
  {
    "address": "chihuahua1nugyyslvtw7fkr7lk5rvffqpy4teurxz4dx3mn",
    "amount": "747454305"
  },
  {
    "address": "chihuahua1nudggsc0v08x75u2vpwq542qx4gfeu7jrzmauz",
    "amount": "197719069"
  },
  {
    "address": "chihuahua15pv2wyqmpramn6sjtdty28t2qs2p5zy7atqamk",
    "amount": "811155156"
  },
  {
    "address": "chihuahua15y6j2pd7prlmrw65tnyuuf69gfwf2dpad2ealv",
    "amount": "106413417"
  },
  {
    "address": "chihuahua159qje5a8u5psws536fhrgt78ppxntuxhkgukzg",
    "amount": "5832461342"
  },
  {
    "address": "chihuahua159rr3ahjkdyada0ung6v2v2yg4d29mlcjrszjg",
    "amount": "1013943"
  },
  {
    "address": "chihuahua15fvungclq4mtx3te8x9c6dhmn4wq4dl532p90w",
    "amount": "247554414"
  },
  {
    "address": "chihuahua15dt6rj00nx2qdv7qdvhe8763jtea5xfl8k76ds",
    "amount": "25455060536"
  },
  {
    "address": "chihuahua15wgfutvyqg3cykgwecg5hdq0apgz35rxq88cl8",
    "amount": "1157457552"
  },
  {
    "address": "chihuahua15cq4djvywu98p3r5a2nm8xss9xryse0krgmda9",
    "amount": "481623373"
  },
  {
    "address": "chihuahua156kuw4y732m923cm7ufq598kkz8tpagyjguvta",
    "amount": "1068691848"
  },
  {
    "address": "chihuahua15ad424dugpac7kud4202s3jax8rgychmxtyyne",
    "amount": "32162015"
  },
  {
    "address": "chihuahua157mzy6th8twrykgh84tgzzjur6r40svk66quuq",
    "amount": "299559599"
  },
  {
    "address": "chihuahua14pghfarja6a9q93ug9hfsuc7w6fs46klnn5h94",
    "amount": "124348418621"
  },
  {
    "address": "chihuahua14fketv99hlrlk80mkggw643spsj3yyf7qylvqn",
    "amount": "1140788434249"
  },
  {
    "address": "chihuahua142zrdg7e8spfgypz7xhkry65mkuw26lthukjta",
    "amount": "33827958358"
  },
  {
    "address": "chihuahua14wrg9c92e5x42u8j85kpqg2yldfx80g6v36qxp",
    "amount": "407337800"
  },
  {
    "address": "chihuahua143cmwm5gkn7rqr3yhm3ajfjvhqupx3eukjam0h",
    "amount": "5069719"
  },
  {
    "address": "chihuahua14khvd82d582zhns8nun04p8qje03zj36z36r0m",
    "amount": "0"
  },
  {
    "address": "chihuahua1kq44na4fs6qfwj6auvz6k0cv2u834t05sk55pd",
    "amount": "6660800564"
  },
  {
    "address": "chihuahua1kptvdkc73ev396rlu059qthy0d8324z0yysgzk",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1kpel43ese325gprc9cmwm2nfqhhzx6tn9xe8hj",
    "amount": "558176649"
  },
  {
    "address": "chihuahua1ksfnaap6jm2ejy9lw268hxnu053xy00uzm5vvk",
    "amount": "674272"
  },
  {
    "address": "chihuahua1knjgaa709qlcrsx5af23ex82dw5yr955ymar5v",
    "amount": "11822586"
  },
  {
    "address": "chihuahua1kkrhn5gucm9ruugz7xrlgjdvapyeg7nza24v4v",
    "amount": "506971"
  },
  {
    "address": "chihuahua1kksufautx29fvmjvm60m9w02uv75a33r95k6s7",
    "amount": "113308235"
  },
  {
    "address": "chihuahua1kmdwjsj3hvpta6v659f7lkhrsge7mq7m8ge7hm",
    "amount": "1576682"
  },
  {
    "address": "chihuahua1hvjylnwyhcz0vu0pdyzt462e9x4lr7tw0z52xl",
    "amount": "15209159"
  },
  {
    "address": "chihuahua1hvhyrhf77s37kcnelzcppqftmz2g52qx0dqxte",
    "amount": "542460010"
  },
  {
    "address": "chihuahua1h50algz9y8nlry39eekpkph2xc689hlxjvm8sf",
    "amount": "19989"
  },
  {
    "address": "chihuahua1h67h9gqyhlgwgtys42qfs26py2gd9mlh4ux7nj",
    "amount": "860219448"
  },
  {
    "address": "chihuahua1cpat972rh4e88xmxlkfk3cg5h5zat7undlc9ju",
    "amount": "425243020"
  },
  {
    "address": "chihuahua1cre6cqr0fuw8uga56d4zz2gxg6m7mjvu455gzh",
    "amount": "253485986"
  },
  {
    "address": "chihuahua1c9mw8a7sytfezamz8u8ftleyyl8dt6tkjet8rg",
    "amount": "13606793713"
  },
  {
    "address": "chihuahua1cg6mvzhuvqhdajjdaksepv2j34rp7xqcrz54vp",
    "amount": "506971"
  },
  {
    "address": "chihuahua1cn7fqru45dejr696jxutys997nlxk6906rlrt3",
    "amount": "6412787"
  },
  {
    "address": "chihuahua1ckj2f8j5dv7n93vlrvmrhmmyjqqj69ph7f45sg",
    "amount": "50697197"
  },
  {
    "address": "chihuahua1cc7py25mkrzsuze5tn9jlf0jju72yr3cw54nfn",
    "amount": "25353668349"
  },
  {
    "address": "chihuahua1evxrund5em6f3dey5xqqcs8pkq9d5306r5rkfx",
    "amount": "253485"
  },
  {
    "address": "chihuahua1esketaa3y435knwqxp5zynp529e3jskvz654jq",
    "amount": "50697197"
  },
  {
    "address": "chihuahua1e5ltjpzwntfxzkvu3antm9kg8dr7wpfj9rlw7e",
    "amount": "114068693"
  },
  {
    "address": "chihuahua1e6294naa8g5jr7vetpcu6m8g9rygwzw4997sft",
    "amount": "15361250"
  },
  {
    "address": "chihuahua1euqrafkcm96mcp0jnuc2qg3fr6vcn3m945qn69",
    "amount": "5069719"
  },
  {
    "address": "chihuahua1mqrjxdsw7tnzq6yzp0287hxsxcmhqyck6xkje6",
    "amount": "207858508"
  },
  {
    "address": "chihuahua1mxma4g0wk3wsftaxdf2y2khy8ecgen9e8s9z7s",
    "amount": "2012678731"
  },
  {
    "address": "chihuahua1m4tqkvw22zgl4w9nyzfeyccnk696r6w7gh3rru",
    "amount": "33797945617"
  },
  {
    "address": "chihuahua1mu89jf6zlc3w3nv4tqeenz20ed50nzup0z84ks",
    "amount": "11693815569"
  },
  {
    "address": "chihuahua1ml0rtgjtk5jv4sj0rlzu79dx2qxa524rg6e8x8",
    "amount": "4664142"
  },
  {
    "address": "chihuahua1uqqrf87u29vhvwy939wj6dh9d28su3f5czn96j",
    "amount": "54166196073"
  },
  {
    "address": "chihuahua1ugmwefslpqpqyh9rrj84fz00fvtd70469y08rx",
    "amount": "25348598"
  },
  {
    "address": "chihuahua1uta2mx7jdxc0tydznryd8a8we6p3xntsvrka26",
    "amount": "50697197"
  },
  {
    "address": "chihuahua1uwkxh7re3ktpckvtsl8jmhylxq7cfwtzvgykwl",
    "amount": "1674841220"
  },
  {
    "address": "chihuahua1ujdhdxtdqmw4e2cey5vaz8836xpqmkq6d4yucr",
    "amount": "1612424358"
  },
  {
    "address": "chihuahua1u5r0j0jytzk2c5z5wswlzy4h385elwwxwumxgg",
    "amount": "36602393"
  },
  {
    "address": "chihuahua1u4ej0v0fwa8r2dd6mkya6vkux3hzg6cjmwlsqu",
    "amount": "15209159"
  },
  {
    "address": "chihuahua1ucf3lyhls8rtxx9wte7e4zhuzvr97g400vrlfy",
    "amount": "32541010005"
  },
  {
    "address": "chihuahua1u7htayu9f9qe4ee4enwtu89ajh0ma29lun04k3",
    "amount": "10207638839"
  },
  {
    "address": "chihuahua1afpxffem0adxjuhdca37pwejhta9spl25dlptq",
    "amount": "10139439"
  },
  {
    "address": "chihuahua1a2zwjvetkqavejdh3nz4cdmccfqq3n6gmz6ugc",
    "amount": "595387884"
  },
  {
    "address": "chihuahua1ak9eeaefhthvu03glvatj05hcmwdpjxsugpw3q",
    "amount": "506971"
  },
  {
    "address": "chihuahua1a7p6zpjnt3sy5wke2lvegw6ltkzktwuthmcphl",
    "amount": "11173713432"
  },
  {
    "address": "chihuahua17rns6zc7rtff3rrzcdy58ggq5qf4cgtrwkqdf9",
    "amount": "1611161669"
  },
  {
    "address": "chihuahua17rl37mp46ntr5dwm7sxtv3ue2zwzcs4lf69ayn",
    "amount": "18250991"
  },
  {
    "address": "chihuahua17ye9r68ty5nr4yvz77z4kd5u3za8v2fanthp33",
    "amount": "10139439"
  },
  {
    "address": "chihuahua1782437t543wlh4ayu5sazsstwxwuy4d4afp7al",
    "amount": "506971"
  },
  {
    "address": "chihuahua17sfxm6wq5tlt2tpnw77n8zh5vkpljxdr7s26tx",
    "amount": "261380817"
  },
  {
    "address": "chihuahua17n6z8ggd4pd2pa0nzp9jspy6fauu69nk8juk66",
    "amount": "2632549657"
  },
  {
    "address": "chihuahua174xlxtk62lnc8xdrenhn3qezt9tvhzphdtzrsv",
    "amount": "567493"
  },
  {
    "address": "chihuahua17439uuzje7tjy3qtmgkzgfpaw9a7vkg50a7p59",
    "amount": "257541762"
  },
  {
    "address": "chihuahua17h628d2wtlw6844nzv0ktnq9qtm5qm3phv0txf",
    "amount": "160608720"
  },
  {
    "address": "chihuahua17eq9z20f6k4qmttw8d2kwfd6vcpxeh0g3glf7d",
    "amount": "11939519816"
  },
  {
    "address": "chihuahua17m4804kq3yd74p090fcm6dg44vp8xq2d67gmv0",
    "amount": "5069719"
  },
  {
    "address": "chihuahua17lamkd8g37ajjm5xe9hchje0ue0v70zraaycrj",
    "amount": "112801263"
  },
  {
    "address": "chihuahua1lpqya4zef3vgxu4ffsk4e0htwr68jcwelddxqq",
    "amount": "506971972"
  },
  {
    "address": "chihuahua1lp5dukv7aq7m5hs6nylyct7xzzkr308lu3ld0k",
    "amount": "1996202142"
  },
  {
    "address": "chihuahua1lp5csm0n8e7a0yq758qhtwwln037wqc6fu4sw7",
    "amount": "8112666899"
  },
  {
    "address": "chihuahua1lyhe9syygu6nwkf3jsg6vs4xneqektltzu5dmy",
    "amount": "507479452"
  },
  {
    "address": "chihuahua1l9vamjetdmzvkr7kfp40pfv77vc5jmkmjhwfx2",
    "amount": "30418318"
  },
  {
    "address": "chihuahua1l22dc2vjgrgwl73ef6rl6zp3uwr007ftvk2j2t",
    "amount": "81426188041"
  },
  {
    "address": "chihuahua1ljd6d8quzfn4zev99unxnz5x89frev3vup780t",
    "amount": "8034922164"
  },
  {
    "address": "chihuahua1lmkt8q9hvx96stsf8a607r8zgzftsqgzkumgqr",
    "amount": "6337149"
  },
  {
    "address": "chihuahua1llmetd7pnfja92ntdjsux38lx904nercvyw3e9",
    "amount": "18789330470"
  }
]`
