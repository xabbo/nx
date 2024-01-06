# nx

A command-line toolkit and Go library for Habbo Hotel.

## Installation

To install the command-line toolkit with Go
```sh
go install github.com/b7c/nx/cmd/nx@latest
```

## Usage

### User
#### Retrieve a user's information

```sh
$ nx user xb7c
             Name │ xb7c
           Status │ Offline
      Last access │ 1 January 2024 3:08:53 am (6 days ago)
        Unique ID │ hhus-c09969403c0b73332345a4b0165ef300
          Created │ 11 November 2023 3:01:19 am
           Figure │ hr-3090-42.hd-180-1.ch-3110-64-1408.lg-275-64.ha-1003-64
            Motto │ (no motto)
  Selected badges │ (none)
```

#### From another hotel
```sh
$ nx user --hotel nl b7c
             Name │ b7c
           Status │ Offline
      Last access │ 12 December 2022 1:58:01 am (1 year ago)
        Unique ID │ hhnl-f4f6928d744d06d4c81aa61116606d25
          Created │ 7 July 2020 6:45:49 am
           Figure │ hr-3090-42.hd-180-1.ch-3110-64-1408.lg-275-64.ha-1003-64
            Motto │ (no motto)
  Selected badges │ (none)
```

You can also set the `HOTEL` environment variable so you don't need to specify `--hotel` every time.
```sh
$ export HOTEL=nl
$ nx user b7c
             Name │ b7c
           Status │ Offline
      Last access │ 12 December 2022 1:58:01 am (1 year ago)
        Unique ID │ hhnl-f4f6928d744d06d4c81aa61116606d25
          Created │ 7 July 2020 6:45:49 am
           Figure │ hr-3090-42.hd-180-1.ch-3110-64-1408.lg-275-64.ha-1003-64
            Motto │ (no motto)
  Selected badges │ (none)
```

#### Outputting the raw JSON response
```sh
$ nx user xb7c --json
{"uniqueId":"hhus-c09969403c0b73332345a4b0165ef300","name":"xb7c","figureString":"hr-3090-42.hd-180-1.ch-3110-64-1408.
lg-275-64.ha-1003-64","motto":"","online":false,"lastAccessTime":"2024-01-03T02:08:53.000+0000","memberSince":"2023-11
-08T14:01:19.000+0000","profileVisible":true,"currentLevel":4,"currentLevelCompletePercent":20,"totalExperience":48,"s
tarGemCount":2,"selectedBadges":[]}
$ nx user xb7c --json | jq .totalExperience
48
```

### Furni
#### Search for furni

```sh
$ nx furni search dragon lamp
Forest Dragon Lamp [rare_dragonlamp*5]
Emerald Dragon Lamp [rare_colourable_dragonlamp*2]
Duck Blue Dragon Lamp [rare_colourable_dragonlamp*5]
Azure Dragon Lamp [rare_colourable_dragonlamp*1]
Bronze Dragon Lamp [rare_dragonlamp*8]
Pink Dragon Lamp [rare_dragonlamp_pink]
Teal Dragon Lamp [rare_colourable_dragonlamp*3]
Rose Gold Dragon Lamp [rare_blackrosegold_dragonlamp]
Bliss Dragon Lamp [nft_ff23_v7_dragon_bliss]
Rainbow Dragon Lamp LTD [rainbow_ltd21_dragonlamp]
Sky Dragon Lamp [rare_dragonlamp*7]
Brown Dragon Lamp [rare_colourable_dragonlamp*4]
Silver Dragon Lamp [rare_dragonlamp*3]
Blue Dragon Lamp [rare_dragonlamp*1]
Diamond Dragon Lamp [diamond_dragon]
Black Dragon Lamp [rare_dragonlamp*4]
Fire Dragon Lamp [rare_dragonlamp*0]
Maroon Dragon Lamp [rare_dragonlamp*10]
```

#### Show furni info

```sh
$ nx furni info 'rare_colourable_dragonlamp*1'
               Name │ Azure Dragon Lamp
        Description │ Scary and scorching!
         Identifier │ rare_colourable_dragonlamp*1
               Type │ Floor
               Kind │ 9136
           Revision │ 69009
               Line │ rare
        Environment │
           Category │ other
  Default direction │ 4
               Size │ 1 x 1
        Part colors │ [#FFFFFF #13ABEA #13ABEA #FFFFFF #FFFFFF #FFFFFF #FFFFFF]
           Offer ID │ -1
             Buyout │ false
                 BC │ false
   Excluded dynamic │ false
      Custom params │
       Special type │ 1
       Can stand on │ false
         Can sit on │ false
         Can lay on │ false
```

### Figure

#### Show figure info
##### By user name
```sh
$ nx figure info -u xb7c
hr-3090-42.hd-180-1.ch-3110-64-1408.lg-275-64.ha-1003-64
┌─ Hair (hr)
│  └─ 3090
├─ Face & Body (hd)
│  └─  180
├─ Shirts (ch)
│  └─ 3110
├─ Trousers (lg)
│  └─  275
└─ Hats (ha)
   └─ 1003
```

##### By figure string
Output all info (`--all`) - parts (`-p`), colors (`-c`) and clothing furni identifiers (`-i`)
```sh
$ nx figure info --all hr-4090-61.hd-180-1.ch-3934-110-110.lg-3596-110-110.ea-3978-110-110.cc-4184-110-110
┌─ Hair (hr)
│  ├─ 4090: Middle Part Hair [clothing_middlepart]
│  │  ├─ hr-4023 [hair_U_middlepart]
│  │  └─ hrb-4023 [hair_U_middlepart]
│  └─   61: #2d2d2d
├─ Face & Body (hd)
│  ├─  180
│  │  ├─ bd-1 [hh_human_body]
│  │  ├─ ey-1 [hh_human_face]
│  │  ├─ fc-1 [hh_human_face]
│  │  ├─ hd-2 [hh_human_body]
│  │  ├─ lh-1 [hh_human_body]
│  │  └─ rh-1 [hh_human_body]
│  └─    1: #ffcb98
├─ Shirts (ch)
│  ├─ 3934: Macho Tattoo [clothing_r20_tattoo]
│  │  ├─ ch-3633 [shirt_M_tattoo]
│  │  ├─ ls-3633 [shirt_M_tattoo]
│  │  ├─ rs-3633 [shirt_M_tattoo]
│  │  ├─ ch-3634 [shirt_M_tattoo]
│  │  ├─ ls-3634
│  │  └─ rs-3634
│  ├─  110: #1e1e1e
│  └─  110: #1e1e1e
├─ Trousers (lg)
│  ├─ 3596: Harem Pants [clothing_harempants]
│  │  ├─ lg-3005 [trousers_U_harempants]
│  │  └─ lg-3006 [trousers_U_harempants]
│  ├─  110: #1e1e1e
│  └─  110: #1e1e1e
├─ Goggles (ea)
│  ├─ 3978: Sleep Time [clothing_r20_slumberoutfit]
│  │  ├─ ea-3720 [acc_eye_U_sleepingmask]
│  │  └─ ea-3721 [acc_eye_U_sleepingmask]
│  ├─  110: #1e1e1e
│  └─  110: #1e1e1e
└─ Jackets (cc)
   ├─ 4184: Kimono by -Push, [clothing_r21_kimono3]
   │  ├─ cc-4218 [jacket_U_kimono3a]
   │  ├─ lc-4218 [jacket_U_kimono3a]
   │  ├─ rc-4218 [jacket_U_kimono3a]
   │  ├─ cc-4219 [jacket_U_kimono3a]
   │  ├─ lc-4219
   │  └─ rc-4219
   ├─  110: #1e1e1e
   └─  110: #1e1e1e
```
