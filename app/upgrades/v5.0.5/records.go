package v505

// Detailed process of how the below list was obtained and amounts calculated to add for the 5% slash and missed days on earning rewards

// This is the list AFTER slashing
// The Silver Fox
// curl -sL 'https://chihuahua-api.polkachu.com/cosmos/staking/v1beta1/validators/chihuahuavaloper124mp8y5prmrumcq9sv6veykf3gutmph5wu2mku/delegations?pagination.limit=10500&pagination.count_total=true' | jq . > chihuahua-tombstone-delegations.json
// Andreisid
// curl -sL 'https://chihuahua-api.polkachu.com/cosmos/staking/v1beta1/validators/chihuahuavaloper16df4qs7tk523gu2qlcvd505cxcf6n2jkvv0nq3/delegations?pagination.limit=10500&pagination.count_total=true' | jq . > andreisid-chihuahua-tombstone-delegations.json

//Based on the current, 22% APR, 5% slash, and possible 9 days of missed rewards:
// The Silver Fox
// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.23/365)*9+1) | floor) | tostring})' chihuahua-tombstone-delegations.json
// Andreisid
// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.23/365)*9+1) | floor) | tostring})' andreisid-chihuahua-tombstone-delegations.json

var recordsJSONString = `[
	  {
    "address": "chihuahua1q2qyq67arqm4j00xd5rsf3sc5ehnx6mwy6eu5d",
    "amount": "248639641"
  },
  {
    "address": "chihuahua1qwqe792ag6p0fg44nk3f2qcredn4lu69p8d03y",
    "amount": "764310137"
  },
  {
    "address": "chihuahua1qaxx33ja7y3azfkevkutr8k8dl958df2442et0",
    "amount": "591384968"
  },
  {
    "address": "chihuahua1pspu8nhc8clxmrw939ujsutycz8k6tss2srqyh",
    "amount": "4883148804"
  },
  {
    "address": "chihuahua1pn5qqjcqxvhnycurcre69y39ywtmpfytan64dy",
    "amount": "276107036"
  },
  {
    "address": "chihuahua1zq0szpt2u8nl8wn4xfvz08vazn6l5arlp9pxqt",
    "amount": "1026424470"
  },
  {
    "address": "chihuahua1zzp2x0lr03c8nqkwx3c7t2dy5f4rnfjx5usxw9",
    "amount": "50157852"
  },
  {
    "address": "chihuahua1zglrnys4y0zrt7807dfxzff6vyldqhgjhzvf0v",
    "amount": "310644301"
  },
  {
    "address": "chihuahua1zv0t2ca3lmgt465dsnpkdtlfl6qr6rn26u6xtw",
    "amount": "3821550685"
  },
  {
    "address": "chihuahua1zmhfe5e85t0ggfcvx7cjeepwpqvtp988k26nds",
    "amount": "47769383"
  },
  {
    "address": "chihuahua1zmac8k24gyr9wjgqz0juyk9h4rlj2rwt2mx23j",
    "amount": "23427299935"
  },
  {
    "address": "chihuahua1yytnp3nglkuywnkfea8h4uwxu7yxr7sppufuny",
    "amount": "3850726020"
  },
  {
    "address": "chihuahua1yxcs29pgx9saerfnavdrx2ylze059x8mxazvll",
    "amount": "5015785"
  },
  {
    "address": "chihuahua1yx72f2eksw5c0adfrcc9xkyv38czn4gw0ytgdg",
    "amount": "28661630"
  },
  {
    "address": "chihuahua1ytamc4wp4n7kel8wmzz5ct0r6sn6uaa6363fs6",
    "amount": "2300903890"
  },
  {
    "address": "chihuahua1yen9j44w4cxt93z0f2c6gcmut8xg24ypspdryl",
    "amount": "955264840"
  },
  {
    "address": "chihuahua19xk0mge8cp0vz2qx0wclfeuww42jq55nyphe8a",
    "amount": "174773843"
  },
  {
    "address": "chihuahua19gm4nf8xm396kfe7m6zfuylte279rg9n7jk49l",
    "amount": "15286202"
  },
  {
    "address": "chihuahua19fc8ya57pajqc66usn2hx9745y4lpdcggfla28",
    "amount": "1732070078"
  },
  {
    "address": "chihuahua19tdquu3q3pphzkpypfu2s0mp4c98am5vatvl3y",
    "amount": "1681482301"
  },
  {
    "address": "chihuahua19epz0cud0ef2kqv8z6ahdc4dh4tkxc07n4nf6d",
    "amount": "47769"
  },
  {
    "address": "chihuahua196wrgjeah0gxskasrq2e9mdyq9dvv4a0hvwnsc",
    "amount": "4332683089"
  },
  {
    "address": "chihuahua197l3uzss7eeral5va35kaj3p2fcz2gz5ayjm7u",
    "amount": "5159093"
  },
  {
    "address": "chihuahua1xp4y9kdjn6ag3fcwvjaj9s88u9gve2yqk4qh8l",
    "amount": "47215831950"
  },
  {
    "address": "chihuahua1xpk4uuwtpnk22a62zp8dsna3ky34ddrrh6mwe0",
    "amount": "8390071221"
  },
  {
    "address": "chihuahua1xr3hycmd438mdrumdflwzdg5666gfd4j36q28k",
    "amount": "519126970"
  },
  {
    "address": "chihuahua1xylkr6vfk8xtfpjwmnsdzc8g3yrzfxhvz6qcx2",
    "amount": "39538002237"
  },
  {
    "address": "chihuahua1xx8d09k74fq8w2j0qxakm30vxn92try6sz4wn7",
    "amount": "1365436646"
  },
  {
    "address": "chihuahua1x2p00vfqlfwntrcfz3d9n74pyzdy808v3x77w3",
    "amount": "2388469178"
  },
  {
    "address": "chihuahua1xjapfntlvpz2uf738lud7yvxcutkvk242vy4jq",
    "amount": "282116919"
  },
  {
    "address": "chihuahua1x7ndut7qw9jtqwgr3ewjwlp4hw8cjckvqnlyhs",
    "amount": "23980230"
  },
  {
    "address": "chihuahua18gkn6pe27nvx653nuz8v40shqal34eh90zprsa",
    "amount": "1969149529"
  },
  {
    "address": "chihuahua18fnz8spzuwswge0k7agz7dn3dn4agsm2h9njte",
    "amount": "577914002"
  },
  {
    "address": "chihuahua18sd2ujv24ual9c9pshtxys6j8knh6xaeatuexm",
    "amount": "15782551"
  },
  {
    "address": "chihuahua18smgsfadpg320pclyn5l6h2n8fx0w0t7qhzgtw",
    "amount": "2736606412"
  },
  {
    "address": "chihuahua1gylzr0qr6ekq6cpwcjlttqtyyveen74grq7zzc",
    "amount": "472406932"
  },
  {
    "address": "chihuahua1g9n76k9h8574r58pyha2dwz8ce3t72slsweyc3",
    "amount": "1924606824379"
  },
  {
    "address": "chihuahua1g2mjkrugdjr0h75cwg9ml5r8ynltd4z9kjy730",
    "amount": "5875634"
  },
  {
    "address": "chihuahua1gtukyuw0fctt3f5f234u46x3pcplas0ncmagun",
    "amount": "2530869710"
  },
  {
    "address": "chihuahua1gnz3pxs5p2t7ealtwy4ehjstdv0sr09mddxcg9",
    "amount": "614852068"
  },
  {
    "address": "chihuahua1g7q5lan32a9qyhcp06gnrdxaf8k5q4g5und6rg",
    "amount": "47769"
  },
  {
    "address": "chihuahua1fxgj0ekxqqe9p2jj8a4q48ann7w7n9zsah4c7n",
    "amount": "366092364"
  },
  {
    "address": "chihuahua1fhngvfsxwxjdd90tyqehm5yjxz85cnq3q0ju3r",
    "amount": "57323260"
  },
  {
    "address": "chihuahua1fl8uhk0q9c4cpzkz82g9mvwhqnwm78crn0gzfj",
    "amount": "339665858"
  },
  {
    "address": "chihuahua128c9zkcmj9886n722qrtdugmgqhsgpny66gv4m",
    "amount": "191077534"
  },
  {
    "address": "chihuahua12vghx05jf5e2sxgfkdhktl3ch0e2k9ae6jc7v9",
    "amount": "1408962375"
  },
  {
    "address": "chihuahua120v78s0lf7umdmwnhct202ke37r9ymnewtsucu",
    "amount": "109821812"
  },
  {
    "address": "chihuahua12saxy20x0x3n7efuxpfeq4jqp0vqpazy8vchhe",
    "amount": "16719332017"
  },
  {
    "address": "chihuahua12s7d79dzts35v30eexrmal0v227xj5nd2zem0z",
    "amount": "23884691783"
  },
  {
    "address": "chihuahua12nel9nty98xfw28p5pez3fks89e0d0rdxglgq6",
    "amount": "2642487419"
  },
  {
    "address": "chihuahua124mp8y5prmrumcq9sv6veykf3gutmph5at2x9q",
    "amount": "307984879551"
  },
  {
    "address": "chihuahua12mvexsd4dfpyy03pjv70wzgwzmm5r5e737v2er",
    "amount": "1882157388"
  },
  {
    "address": "chihuahua1tvcd8welups8wwgj9jxqpt9027eyzatlk5ytpl",
    "amount": "10986958220"
  },
  {
    "address": "chihuahua1td069hed3jwkv6c38nyrcmmpdz24r3cfyl60zz",
    "amount": "509893415"
  },
  {
    "address": "chihuahua1t5gyhdy29snzaxzsn63fr7pezn203emevz87f9",
    "amount": "493480001"
  },
  {
    "address": "chihuahua1vqwx49ymv2czkuk4vrtdy0rkz8laec9lcc6c76",
    "amount": "31854370949"
  },
  {
    "address": "chihuahua1vgj4cqrk9jt9m743d8skqpmh08jatym5wrehca",
    "amount": "2355984593"
  },
  {
    "address": "chihuahua1vwpel70r8rurhs5fu55v9us5775fe93mdnskq0",
    "amount": "532578004"
  },
  {
    "address": "chihuahua1vs544qt0e24mhh6nxkf4d7zwudqqy79axerhu3",
    "amount": "113070130"
  },
  {
    "address": "chihuahua1dg459uca2ga7rufus3ykwl52tljdtagsedv3ht",
    "amount": "145456881"
  },
  {
    "address": "chihuahua1dcvlnewq7ktr5kd47755d4x7jvzsn3ae2g2jj8",
    "amount": "8598489"
  },
  {
    "address": "chihuahua1d6486j6gmapqexvxpd4wesah5c0w4vvjq52rkn",
    "amount": "104376103"
  },
  {
    "address": "chihuahua1w0w9rhf6fjd3rd3lsu9t02nteg5uxksjlkzhvk",
    "amount": "11178035"
  },
  {
    "address": "chihuahua10z6mwcv78kmnj4lm74ce3zydrqwuh5qhu36pdq",
    "amount": "3479521898"
  },
  {
    "address": "chihuahua109jutq7jyxqycls4fpn9shp9ep8370vyz0ksv4",
    "amount": "1433081507"
  },
  {
    "address": "chihuahua10dutcucqly40qr50a2nzq0ej3zjq2u83g24z30",
    "amount": "23884691"
  },
  {
    "address": "chihuahua10sy9drt0znfc4xcck770ptx88zqnc4rqne4jzj",
    "amount": "716540753"
  },
  {
    "address": "chihuahua10e6lxuc6kf0a2xs2n6e8tu9klx5jaee6y4r77u",
    "amount": "1307304720"
  },
  {
    "address": "chihuahua1sqfzrysfghnpljm30q6qw3dv4z5x5c0qhp4rj3",
    "amount": "391648951"
  },
  {
    "address": "chihuahua1spq8qzxns9asz64azsu2pt37pv8dqj4r3lxcej",
    "amount": "11942345891"
  },
  {
    "address": "chihuahua1s4guk3hpgqk4nycwa49rdz8qxh424catms732e",
    "amount": "981913"
  },
  {
    "address": "chihuahua13zs0yntzdswanglzutk89zad74dwc4v85qx5cd",
    "amount": "191077534267"
  },
  {
    "address": "chihuahua13vcc39q4v8aupx0ccxxhscg0l9karkjwgvc03c",
    "amount": "1929883096102"
  },
  {
    "address": "chihuahua1332gmq88rwjgj2jz5awhp5v4duu6ngukcfsghr",
    "amount": "122098544"
  },
  {
    "address": "chihuahua13jawsn574rf3f0u5rhu7e8n6sayx5gkw3eddhp",
    "amount": "57017384"
  },
  {
    "address": "chihuahua1jz9w0ewc0u4wmtvjggmpckpq8qxzz2jmtj8ppd",
    "amount": "518408378"
  },
  {
    "address": "chihuahua1jzdhw6dnw6jzaamh5njtm9kmdf0t46xnhp90mn",
    "amount": "3817920212"
  },
  {
    "address": "chihuahua1jyycnfkafcgjap9hx4hwwz96zqwv0g63cr2aa9",
    "amount": "997232720"
  },
  {
    "address": "chihuahua1j800ng2kgwqtxsfw2xnrujau0xfxaa5u37pdmw",
    "amount": "9601646"
  },
  {
    "address": "chihuahua15gn6xpxr7psf4vcaspkf8fg4axtyvw20sfjdy2",
    "amount": "313238353"
  },
  {
    "address": "chihuahua15e5a5p5qp46e8wx3h3ypetdrqe822pqpqzjrqh",
    "amount": "1385312"
  },
  {
    "address": "chihuahua15aq8hur3uz7yz2q8juewm2s2rjchgma0m5zv3n",
    "amount": "461605"
  },
  {
    "address": "chihuahua14rh69nvuz0m7j5lmutckevrle6z5wvfk7ddlsa",
    "amount": "3893552399"
  },
  {
    "address": "chihuahua14xc4l7f3c3gde0r0xzm5qm7azpd4sgyrkqdfj2",
    "amount": "66877136"
  },
  {
    "address": "chihuahua14gc2aedc7tm6mn9zgp77g9dfunqpqzq0zzx7hl",
    "amount": "30307122"
  },
  {
    "address": "chihuahua14fketv99hlrlk80mkggw643spsj3yyf7qylvqn",
    "amount": "3611365397656"
  },
  {
    "address": "chihuahua1438svnqnj4t3h0wmrtj007mjakgh0cmasr2uap",
    "amount": "58797103"
  },
  {
    "address": "chihuahua1442dcsssrha9798205nunsam06da2grdlsv6kt",
    "amount": "1170599340"
  },
  {
    "address": "chihuahua1kytdumzek99q6x860l4lru6rtvc68h58c6lxsj",
    "amount": "1499958643"
  },
  {
    "address": "chihuahua1k8wnqfukgm55fcfzsnpqf9z8a3s89jsgljlalx",
    "amount": "142973765"
  },
  {
    "address": "chihuahua1k33zpxwjc8chrshrldrm4j0a8puyv86rekacry",
    "amount": "2299865434"
  },
  {
    "address": "chihuahua1k53z9k72kn4h65wfjaq8k2p9kcf4ezjhuzgs2c",
    "amount": "362814997"
  },
  {
    "address": "chihuahua1hr6lckc4nyyjdmgfadyq7drep7vxmzqvsccjma",
    "amount": "2398877"
  },
  {
    "address": "chihuahua1hy72n6amu2apz04tw730d8smm3fg90ygvd4rmw",
    "amount": "22126300774"
  },
  {
    "address": "chihuahua1cre6cqr0fuw8uga56d4zz2gxg6m7mjvu455gzh",
    "amount": "238846917"
  },
  {
    "address": "chihuahua1cn7fqru45dejr696jxutys997nlxk6906rlrt3",
    "amount": "6538676"
  },
  {
    "address": "chihuahua1c7j5k22ufmnyh8fw5sr8r9en7dg7uk64cr07xj",
    "amount": "7289290"
  },
  {
    "address": "chihuahua1cl972hflpyh59fztds6m0h3n6879cmhgry89vl",
    "amount": "14330815"
  },
  {
    "address": "chihuahua1e8gddqsrku6j9px3qwsp7kd0lruastrh84yk3r",
    "amount": "14720871692"
  },
  {
    "address": "chihuahua1evxrund5em6f3dey5xqqcs8pkq9d5306r5rkfx",
    "amount": "238846"
  },
  {
    "address": "chihuahua1e6294naa8g5jr7vetpcu6m8g9rygwzw4997sft",
    "amount": "47769383"
  },
  {
    "address": "chihuahua1e6ltmckcj42cdp6yx07qu8k4hxsucgcfxrdv2u",
    "amount": "20540834"
  },
  {
    "address": "chihuahua16crst3vutxuvaksay8wa4q3twll4ss4xs54k63",
    "amount": "1025508194"
  },
  {
    "address": "chihuahua1m8gs2emktxx53tfunx6qfleesymy2aje7jxwrc",
    "amount": "31050099318"
  },
  {
    "address": "chihuahua1mker6v4vu3eecetfx43d7gntm95d4hsadvdl7v",
    "amount": "448602281"
  },
  {
    "address": "chihuahua1mu26eqsel6ft7zf4ce596nwcmp5umegj40l0dm",
    "amount": "1310783940806"
  },
  {
    "address": "chihuahua1masrhdyavztuz5f7gc08ga3vajvjjkh49k0l3v",
    "amount": "257845085"
  },
  {
    "address": "chihuahua1ufppyphwfwlhxft05fuj0206kq9fheslmz5p2m",
    "amount": "10095820369"
  },
  {
    "address": "chihuahua1u05gu68e9xkp2jhgq2792arkqcla7dwvm8k7tk",
    "amount": "14321722"
  },
  {
    "address": "chihuahua1a8gx9m553deeajj8jcdlp37tgwqyy2c00nsjxh",
    "amount": "87876868510"
  },
  {
    "address": "chihuahua1a0hredawn0h50e2lhvr6rmmyjshkjpjkgeeekw",
    "amount": "1413136"
  },
  {
    "address": "chihuahua1amxdm09q53wtzf7cm3470ady3urn96gmgf5cfh",
    "amount": "1433081"
  },
  {
    "address": "chihuahua178a0lrsumcrp7f0dmkm400h87qgwxa3c0gjule",
    "amount": "11990115"
  },
  {
    "address": "chihuahua17v20kpuelhgpz2lqxhvhnw9g95qurdt5mkmmyr",
    "amount": "104949335"
  },
  {
    "address": "chihuahua174xlxtk62lnc8xdrenhn3qezt9tvhzphdtzrsv",
    "amount": "583241"
  },
  {
    "address": "chihuahua1lx20fvhscnaeqlypqr5wwuajdemvwup7vt89yw",
    "amount": "3511464247"
  },
  {
    "address": "chihuahua1l238kq8vc2n233lmcxc5urq35yqc5zsa7ahmn4",
    "amount": "324354114"
  },
  {
    "address": "chihuahua1ljd6d8quzfn4zev99unxnz5x89frev3vup780t",
    "amount": "25026618897"
  },
  {
    "address": "chihuahua1qg7nwmxs2g0fp2t7fp5csgytk2gxdlh0cs35c6",
    "amount": "178476626"
  },
  {
    "address": "chihuahua1qnpx3y9jrx6z439r55aekmelkwe3fp42dte359",
    "amount": "477693835"
  },
  {
    "address": "chihuahua1qly96sy3yklkwv6cf4he5ves4qu8yn58ea7wn5",
    "amount": "48102067"
  },
  {
    "address": "chihuahua1ql3zye37yaczg8x9tvg8uutnwgj88u3qedr8ac",
    "amount": "138531212"
  },
  {
    "address": "chihuahua1py7vnejer4auv8e0dyexvgkas4q08wppjyl99y",
    "amount": "95538767"
  },
  {
    "address": "chihuahua1p2yguyfencsn3zwtxhqdukf9m0pcqfy3nepnhn",
    "amount": "238846917"
  },
  {
    "address": "chihuahua1pvf4ju252as6yzl6lasc45a348vhcnz5u8053l",
    "amount": "6233904556"
  },
  {
    "address": "chihuahua1p4dux668yurw3zzvuw6apg25hf3fxnayrmnlcu",
    "amount": "1153821690"
  },
  {
    "address": "chihuahua1zduvf6awl2rcjdtq58g69cejf2d8unkqevsx7f",
    "amount": "59663960"
  },
  {
    "address": "chihuahua1zwt7pkrkuz0zasa4mjge6we67mpvngxe2mtpx7",
    "amount": "3483747644"
  },
  {
    "address": "chihuahua1z6a9d2ta4dtedemzt6vzlhphmyv6fytx6jk94m",
    "amount": "39648588"
  },
  {
    "address": "chihuahua1ryu2zpqjyf6zd6cxhlcszrg4ykt5lk2arvc962",
    "amount": "46903945652"
  },
  {
    "address": "chihuahua1rgnnvltmwujr2yf067c6n8pzwz99fquznsz9q7",
    "amount": "1194234"
  },
  {
    "address": "chihuahua1rfk6ermke7cfvxz5ruyt6e3mf3gjpztnmduxx3",
    "amount": "526418606"
  },
  {
    "address": "chihuahua1rwrsx5fel5zfqdt6a5lelrqrhsk3r3s9pfk479",
    "amount": "4776938"
  },
  {
    "address": "chihuahua1rlrdcndumuvrtl55xttvdsy0qfdes7ua0wzk6m",
    "amount": "191077534"
  },
  {
    "address": "chihuahua1yr73cp3unx54rek5gr4twyr8twj3wgewt80c09",
    "amount": "1872512066"
  },
  {
    "address": "chihuahua1yur5ph8x32sa5dx8akz6a9cljvl6fvr9wfrmsu",
    "amount": "53262862"
  },
  {
    "address": "chihuahua19p48jkks89cm2u46ea993wuxglnflrtaqnsv93",
    "amount": "270804635"
  },
  {
    "address": "chihuahua19kvla8xpm44zp0gj0gdfqfggh5m60p8r9p37ug",
    "amount": "29712556"
  },
  {
    "address": "chihuahua19ckpdvfayewdgrttw782qhv70fy6l2wj55k4hz",
    "amount": "90761828"
  },
  {
    "address": "chihuahua196mwqvttaczyh6mygd9k9shwts70u227880mwr",
    "amount": "3407199052"
  },
  {
    "address": "chihuahua197l3uzss7eeral5va35kaj3p2fcz2gz5ayjm7u",
    "amount": "621001"
  },
  {
    "address": "chihuahua1xxpd2e9apcw3l55ku46d09jmdp3f2h9grfqc2t",
    "amount": "26846393"
  },
  {
    "address": "chihuahua1xg3jhs6d0s56gkrzwqjrc26tuafxxu6wrwvhy3",
    "amount": "47821934"
  },
  {
    "address": "chihuahua1xsf5rmhg62740lwh5ntyzn0v69q47m9hcfw206",
    "amount": "726505542"
  },
  {
    "address": "chihuahua1xjapfntlvpz2uf738lud7yvxcutkvk242vy4jq",
    "amount": "177121782"
  },
  {
    "address": "chihuahua1xnpfjs9cu4fp3hx8ur3a4xgzfzk74pg3alvgrn",
    "amount": "25795467"
  },
  {
    "address": "chihuahua1xeczmgzhhxg0fghwfw8602zhhg46nzvqpjfp53",
    "amount": "23884691"
  },
  {
    "address": "chihuahua1x69jls25lldyytfhzaaved078v72rc3js3acmt",
    "amount": "2560056804"
  },
  {
    "address": "chihuahua1x6t5t62hn272zz6w2ydjqtau7a44eam5jmcwsk",
    "amount": "84697505"
  },
  {
    "address": "chihuahua1xmrhxyk8h2tuntnsdwx0e4twpqnmkpn4qm99jv",
    "amount": "77052015"
  },
  {
    "address": "chihuahua1xu53xhwy7x53dqw904flsc5c8wdm7n5gl7u7ph",
    "amount": "5015785"
  },
  {
    "address": "chihuahua189e7ewzw5nxm0tsx286q3aqpwgj9wxkzuyvwr5",
    "amount": "47806424828"
  },
  {
    "address": "chihuahua18sd2ujv24ual9c9pshtxys6j8knh6xaeatuexm",
    "amount": "23258704"
  },
  {
    "address": "chihuahua18h28l77fd86ldtdsa096npmkqrs3w0ruxj0u6j",
    "amount": "238846917"
  },
  {
    "address": "chihuahua1gylzr0qr6ekq6cpwcjlttqtyyveen74grq7zzc",
    "amount": "513440067"
  },
  {
    "address": "chihuahua1gwywmna0jky46mcj3wvc336qgyk9u2n3mtcnru",
    "amount": "47769383575"
  },
  {
    "address": "chihuahua1fqkyvd5lnadfffs4ky7jcnqkcsl2pyaylujv5e",
    "amount": "111456782"
  },
  {
    "address": "chihuahua1fwrdyu089snfj464cwpant96fuzf4xpkeutpkc",
    "amount": "2296165176"
  },
  {
    "address": "chihuahua1fcj7ww85n35tgdamasfye6w6xw2f0zldm49j0f",
    "amount": "7587710383"
  },
  {
    "address": "chihuahua1furgr24vn5dme5fwg4xfu08v44kjd22rr2anp7",
    "amount": "73994775"
  },
  {
    "address": "chihuahua1fl3zm62mslvkfpa0uy2lz6jxe468443w050ag2",
    "amount": "2144845322"
  },
  {
    "address": "chihuahua12ftk2gst84zqn9ny2hcuyntuz4jfnxpnlaz77t",
    "amount": "82315103080"
  },
  {
    "address": "chihuahua12vzkxeu3s4w4gj8s4yxyztnlhsfzf52h9hcx5q",
    "amount": "23884691"
  },
  {
    "address": "chihuahua12s3pvtsk664sa50kw7te2rh735c026l66d6tax",
    "amount": "2218410173"
  },
  {
    "address": "chihuahua1253cfxjzyjw8rkyxq9fnjlzludcvepxtrpanqf",
    "amount": "56606719"
  },
  {
    "address": "chihuahua12mlkq3ruj0ql27cc0ltd54vzxrsjwkkmnssu83",
    "amount": "43905555"
  },
  {
    "address": "chihuahua1tqlex8qlnk6pwphvk95guvucgh994ysehaj8zf",
    "amount": "37878413"
  },
  {
    "address": "chihuahua1thu6fuhdywn5pmw0cm03plz3gwrrlpew28u03h",
    "amount": "244960825740"
  },
  {
    "address": "chihuahua1vyzmrltdxlc5x6et8v592nl8uqlwl7qrycvapd",
    "amount": "23884691"
  },
  {
    "address": "chihuahua1v9fjgctqnj6qd73xspyllr6z5srs7qle8eqq86",
    "amount": "37451"
  },
  {
    "address": "chihuahua1v80n5t3ucpj8dhdve0l80dz77mlk6l44sqe4g5",
    "amount": "716540753"
  },
  {
    "address": "chihuahua1vnqvus65hx399kep83ukhf4myf5t5gm94nz9q5",
    "amount": "4776938"
  },
  {
    "address": "chihuahua1vntrgjjt92vddykr5n30lhcmzvf67kuar7zd56",
    "amount": "214484532"
  },
  {
    "address": "chihuahua1v5rfq8w959phlam7nmmwq7s9k3uwfxl6jk6h0f",
    "amount": "11944232"
  },
  {
    "address": "chihuahua1ve4pzf8mqgzrfkhkg9y4kcn3lle64c49r6u9qe",
    "amount": "382155068"
  },
  {
    "address": "chihuahua1va3nl4wup8zexwysug2pkuu4htgcn78qgggd04",
    "amount": "481276539"
  },
  {
    "address": "chihuahua1d9p6d3njrnaups7sc2twtwnlk8xjck32qnw7sf",
    "amount": "6210019"
  },
  {
    "address": "chihuahua1d0egu5ge0n6m8yk8fm2ev37xhrfqel7qgupevw",
    "amount": "23884691"
  },
  {
    "address": "chihuahua1dh03457vd2e8f4fykdzuetpuxtrxjmpjuqd30e",
    "amount": "30433984"
  },
  {
    "address": "chihuahua1dummey3rjqdydggux6dtnawy2k6eup95ue9qn4",
    "amount": "5354995668"
  },
  {
    "address": "chihuahua1w9n2m9lxfpedg50r0wj7f8wgnv5720s6807epp",
    "amount": "4829489"
  },
  {
    "address": "chihuahua1wfwqp447e2kl5hz5nyvzgurrm02tyrfv6kld3n",
    "amount": "238846"
  },
  {
    "address": "chihuahua1w56pgj76ynzur9xklwy3ge9f8cwl6y6gzfa46m",
    "amount": "9553876714"
  },
  {
    "address": "chihuahua1w4tv7u3kzusl0ze5e5zacn6xnadrgf2mzmzhjx",
    "amount": "3305306957"
  },
  {
    "address": "chihuahua1w4459hvuq5e8ggyrxxst3wggndzgppn0cqpq8x",
    "amount": "28661630"
  },
  {
    "address": "chihuahua10fdwe698fl5jfuq38lh6tgs0f0jp9uak53pkf4",
    "amount": "19107753"
  },
  {
    "address": "chihuahua10uwy4lr7tlax4tmaa96hpmy3s63504qmnu0jw4",
    "amount": "238846"
  },
  {
    "address": "chihuahua10u7k3dupmp5f3x3pfkkjnc6hjxnty2p23jyykx",
    "amount": "267508547"
  },
  {
    "address": "chihuahua1sr0s7t5vskgkdeuqw47zatvdlhfhvuwkv3kv6y",
    "amount": "655491481"
  },
  {
    "address": "chihuahua1sghqztvgjv3ug6ltpvvwwyhpt35u7hu8nz7r7m",
    "amount": "1365214508"
  },
  {
    "address": "chihuahua1st0jrcm9zl7qww4uyjtyhn4trddcafqrzv07f9",
    "amount": "143308150"
  },
  {
    "address": "chihuahua1sevq5s6rdvtkeex4lqvuzvwu7gwjlykutv7ncs",
    "amount": "9553"
  },
  {
    "address": "chihuahua1387zfawsz0gd4vp0jwparfem3sjcfy5aj0pl37",
    "amount": "668771370"
  },
  {
    "address": "chihuahua13vykv68dgcs4d0amu7lm3vegnk3qgaxhjes9ma",
    "amount": "238846917"
  },
  {
    "address": "chihuahua13vcc39q4v8aupx0ccxxhscg0l9karkjwgvc03c",
    "amount": "1129745921552"
  },
  {
    "address": "chihuahua13306hrp8e9y9z8lv5h2eccupffj8kw95uqc5j6",
    "amount": "148615329"
  },
  {
    "address": "chihuahua13jawsn574rf3f0u5rhu7e8n6sayx5gkw3eddhp",
    "amount": "41733155"
  },
  {
    "address": "chihuahua13a04ryy0khxgxtnp9vgwqevqgsw0fxxczhp47z",
    "amount": "46622583"
  },
  {
    "address": "chihuahua1jpv8lhrkaqw5mznvx8yceh42p9xxaahqeq902t",
    "amount": "26215837706"
  },
  {
    "address": "chihuahua1ju3z4nef3ppxz8qed0cqcyytw0s6js5z7j8rw7",
    "amount": "11553980805"
  },
  {
    "address": "chihuahua1nrh06xmlsh2s3x94kdd0mxesufu4mnwywykcva",
    "amount": "12133232"
  },
  {
    "address": "chihuahua1n9035dtgz6wtnk7mfpps4mhvtetfd20yfr7e2j",
    "amount": "5700137363"
  },
  {
    "address": "chihuahua1n4tdkuvyxtgvahwayg4533vzzprarkwulr4apq",
    "amount": "238846917"
  },
  {
    "address": "chihuahua1nc9ffllzatejq8zam0crk058e2dzrjdkt70cx6",
    "amount": "66877136"
  },
  {
    "address": "chihuahua1n7akar93zcs8l9hgwqehz4q5uax8xcry5fhhv7",
    "amount": "1194955321"
  },
  {
    "address": "chihuahua15ddgx30w2l50dvs0vj0f9jxf9f5gfv8x7vgy52",
    "amount": "2734844979"
  },
  {
    "address": "chihuahua153j5kshmp670frzq96xsu90whn6ww8eammz999",
    "amount": "191077"
  },
  {
    "address": "chihuahua15mxxewmh0tdsqjppp0zg6a69erz5vkv88gemqg",
    "amount": "1889183581"
  },
  {
    "address": "chihuahua15aq8hur3uz7yz2q8juewm2s2rjchgma0m5zv3n",
    "amount": "871552"
  },
  {
    "address": "chihuahua15asahu04jrhlledhwmvwtcte6vrz0atr8l7uu7",
    "amount": "1569988560"
  },
  {
    "address": "chihuahua14xaxn5yr7p2ganuastmtvw54wta928ny5cgdqp",
    "amount": "472869127"
  },
  {
    "address": "chihuahua14fp7zfehp25gr78ntn9uv7qmqgx74dkl8gtpf8",
    "amount": "8671490885"
  },
  {
    "address": "chihuahua14fketv99hlrlk80mkggw643spsj3yyf7qylvqn",
    "amount": "1521454866868"
  },
  {
    "address": "chihuahua14h8fznw7tna3fyx4lhrvzvkac28lyfk008kgrs",
    "amount": "31615783249"
  },
  {
    "address": "chihuahua1krqycrpq8n0qwkufhe3ug6d0emxxytledp9puy",
    "amount": "948604419"
  },
  {
    "address": "chihuahua1kuqqvtgv55zl6h4lutd058l63t9nlp07p3x00m",
    "amount": "200631"
  },
  {
    "address": "chihuahua1hxptnsft8ylfsyydpfxlswh5tajwjqu26ldczg",
    "amount": "188712949"
  },
  {
    "address": "chihuahua1hvjylnwyhcz0vu0pdyzt462e9x4lr7tw0z52xl",
    "amount": "6544405"
  },
  {
    "address": "chihuahua1hh45avpvtjurwur6wemr79wtrmf3xcsxz2yzg2",
    "amount": "5184267891"
  },
  {
    "address": "chihuahua1hljmqq8q66c3c0zhaqxm09lc40g35zum75nv9v",
    "amount": "310500993"
  },
  {
    "address": "chihuahua1cf9ad68jrhq569melk4zrqak8r3y9z32dpqufp",
    "amount": "48366760699"
  },
  {
    "address": "chihuahua1cftk8sg8tev26uyvvhg0wmugh7xkazymh47lhp",
    "amount": "23884691"
  },
  {
    "address": "chihuahua1cdak665z2vewlkqwf5z8024qg2flv09yy7hj9s",
    "amount": "189214528"
  },
  {
    "address": "chihuahua1cn4flnq4chrh8tall7fwjvhjjkh0k5k2w2anfr",
    "amount": "17196978"
  },
  {
    "address": "chihuahua1cn7fqru45dejr696jxutys997nlxk6906rlrt3",
    "amount": "7105563"
  },
  {
    "address": "chihuahua1evxrund5em6f3dey5xqqcs8pkq9d5306r5rkfx",
    "amount": "238846"
  },
  {
    "address": "chihuahua1echtjakrrmy7zrjzk727vu5mgevswmr53m58m5",
    "amount": "1292161825"
  },
  {
    "address": "chihuahua1e6294naa8g5jr7vetpcu6m8g9rygwzw4997sft",
    "amount": "14330815"
  },
  {
    "address": "chihuahua1eul3948ahtp9und9d39cjl9jlwtlfrl03ujt7v",
    "amount": "2160227064"
  },
  {
    "address": "chihuahua16psdhg5ty795fhtpwt3kzgcnevwmn65euuq4zt",
    "amount": "463363020"
  },
  {
    "address": "chihuahua16z48nqpsjtpwyvdem0puln0cl2xqd4nwsmhmvz",
    "amount": "47769"
  },
  {
    "address": "chihuahua16r830xe7k54s7kcrheq83uwdfuxvzxnumt8hzz",
    "amount": "1844471438"
  },
  {
    "address": "chihuahua16df4qs7tk523gu2qlcvd505cxcf6n2jklm0wnd",
    "amount": "477693"
  },
  {
    "address": "chihuahua165yw8g8cyy25kptfl63cdvkkcm2v7zwxufr8j4",
    "amount": "76609795"
  },
  {
    "address": "chihuahua16unlqgz8fnq7rynr6g5svygxkqaepghnwgvucf",
    "amount": "30572405"
  },
  {
    "address": "chihuahua1mps9fertlaezrcz24gf9xzl9jwzc2ym6wyratf",
    "amount": "47769383"
  },
  {
    "address": "chihuahua1m9ulj0wd7jvhh0hrg4cw83vf39edwsw64v55jv",
    "amount": "77669584550"
  },
  {
    "address": "chihuahua1mta9uggk77k57xfh8rd0kemzy9ggefygy3xedp",
    "amount": "3236083690"
  },
  {
    "address": "chihuahua1md9zaw7gpnrfl2swql6hjvgttjuf6cvkqlsuwy",
    "amount": "242317"
  },
  {
    "address": "chihuahua1m42gk4aepydq5w7nhawxhgjmv2rqc3huavnwwj",
    "amount": "42992445"
  },
  {
    "address": "chihuahua1mmfl8uz5ff008rr7v76nts8wunsw0nu2s58kyq",
    "amount": "143355920"
  },
  {
    "address": "chihuahua1u2zjcdl3606gz4wg0yc6mnkpudyfv4twnf5tvt",
    "amount": "24337461"
  },
  {
    "address": "chihuahua1uwg60fhx8y5u3x4c7ymqpgjqp5egzqcjhcd6jh",
    "amount": "23884691"
  },
  {
    "address": "chihuahua1u05gu68e9xkp2jhgq2792arkqcla7dwvm8k7tk",
    "amount": "16408251"
  },
  {
    "address": "chihuahua1usl6upygejydag3xj37sl0edln6stctvmhu66f",
    "amount": "1587472154"
  },
  {
    "address": "chihuahua1u5r0j0jytzk2c5z5wswlzy4h385elwwxwumxgg",
    "amount": "4955932"
  },
  {
    "address": "chihuahua1u5w88m7glt08pea0c8xxd0gggpnymn9wgm4t0k",
    "amount": "687629491661"
  },
  {
    "address": "chihuahua1uhc2jpk7z4mlgt9u7xpf2vgleupf3ehwyhsqna",
    "amount": "527985774"
  },
  {
    "address": "chihuahua1uc73gnaj5uyehgzcudcd6h83u324spkvgxh22r",
    "amount": "500164553"
  },
  {
    "address": "chihuahua1ulqdpwj8zgp9aj98uefp7r80qhk6ereunl9us4",
    "amount": "124627589"
  },
  {
    "address": "chihuahua1a2zwjvetkqavejdh3nz4cdmccfqq3n6gmz6ugc",
    "amount": "170536699"
  },
  {
    "address": "chihuahua1a0hredawn0h50e2lhvr6rmmyjshkjpjkgeeekw",
    "amount": "2326501"
  },
  {
    "address": "chihuahua1acxrhexvf7kxk7a7qepgy6ajruwpc6kg2et3sc",
    "amount": "477"
  },
  {
    "address": "chihuahua1a6ah5gpma5rhasyvr3krzlw82v06jvys4krvez",
    "amount": "119423458"
  },
  {
    "address": "chihuahua1aa32ps7czfd3s8fukcauv2tvj9ga83pt5wxvv7",
    "amount": "4192145563"
  },
  {
    "address": "chihuahua179pxpp97hrsmajwcd3jn2t46ekjpnjah62r8dw",
    "amount": "1289773356"
  },
  {
    "address": "chihuahua17crvfnhtc26k3tpwrkgc7v4lj0zf9klsgqe5at",
    "amount": "59329574"
  },
  {
    "address": "chihuahua17m4804kq3yd74p090fcm6dg44vp8xq2d67gmv0",
    "amount": "164162980"
  },
  {
    "address": "chihuahua1lg8ukq2ehc9k0wgjfk0afm7uea8750yps09u74",
    "amount": "47769383"
  },
  {
    "address": "chihuahua1lcyv3t873t2ntv0pljpt3q88fes042765vmwww",
    "amount": "189739991"
  },
  {
    "address": "chihuahua1lmkt8q9hvx96stsf8a607r8zgzftsqgzkumgqr",
    "amount": "5254632"
  }
]`