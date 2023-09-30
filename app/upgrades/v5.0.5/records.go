package v505

// Detailed process of how the below list was obtained and amounts calculated to add for the 5% slash and missed days on earning rewards

// This is the list AFTER slashing
// curl -sL 'https://chihuahua-api.polkachu.com/cosmos/staking/v1beta1/validators/chihuahuavaloper124mp8y5prmrumcq9sv6veykf3gutmph5wu2mku/delegations?pagination.limit=10500&pagination.count_total=true' | jq . > chihuahua-tombstone-delegations.json

//Based on the current, 22% APR, 5% slash, and possible 9 days of missed rewards:
// jq '.delegation_responses[] | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.23/365)*9+1) | floor) | tostring})' chihuahua-tombstone-delegations.json

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
  }

]`