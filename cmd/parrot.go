package cmd

import (
	"flag"
	"fmt"
	isatty "github.com/mattn/go-isatty"
	"github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ParrotCmd represents the init command
func ParrotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "parrot",
		Run:    parrotRun,
		Args:   cobra.NoArgs,
		Hidden: true,
	}

	return cmd
}

func parrotRun(cmd *cobra.Command, args []string) {
	loops := flag.Int("loops", 0, "number of times to loop (default: infinite)")
	delay := flag.Int("delay", 75, "frame delay in ms")
	orientation := flag.String("orientation", "regular", "regular or aussie")
	flag.Parse()

	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		fmt.Fprintf(os.Stderr, "%s must be run in a terminal!\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	termbox.SetOutputMode(termbox.Output256)

	loop_index := 0
	draw(*orientation)

loop:
	for {
		select {
		case ev := <-event_queue:
			if (ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC)) || ev.Type == termbox.EventInterrupt {
				break loop
			}
		default:
			loop_index++
			if *loops > 0 && (loop_index/9) >= *loops {
				break loop
			}
			draw(*orientation)
			time.Sleep(time.Duration(*delay) * time.Millisecond)
		}
	}
}

var frame_index = 0

func reverse(lines []string) []string {
	newLines := make([]string, len(lines))
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		newLines[i], newLines[j] = lines[j], lines[i]
	}
	return newLines
}

func draw(orientation string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	lines := strings.Split(frames[frame_index], "\n")

	if orientation == "aussie" {
		lines = reverse(lines)
	}

	for x, line := range lines {
		for y, cell := range line {
			termbox.SetCell(y, x, cell, colors[frame_index], termbox.ColorDefault)
		}
	}

	termbox.Flush()
	frame_index++
	if frame_index == num_frames {
		frame_index = 0
	}
}

const num_frames = 10

var frames = [num_frames]string{
	`                        .cccc;;cc;';c.
                      .,:dkdc:;;:c:,:d:.
                     .loc'.,cc::::::,..,:.
                   .cl;....;dkdccc::,...c;
                  .c:,';:'..ckc',;::;....;c.
                .c:'.,dkkoc:ok:;llllc,,c,';:.
               .;c,';okkkkkkkk:,lllll,:kd;.;:,.
               co..:kkkkkkkkkk:;llllc':kkc..oNc
             .cl;.,okkkkkkkkkkc,:cll;,okkc'.cO;
             ;k:..ckkkkkkkkkkkl..,;,.;xkko:',l'
            .,...';dkkkkkkkkkkd;.....ckkkl'.cO;
         .,,:,.;oo:ckkkkkkkkkkkdoc;;cdkkkc..cd,
      .cclo;,ccdkkl;llccdkkkkkkkkkkkkkkkd,.c;
     .lol:;;okkkkkxooc::loodkkkkkkkkkkkko'.oc
   .c:'..lkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkd,.oc
  .lo;,ccdkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkd,.c;
,dx:..;lllllllllllllllllllllllllllllllloc'...
cNO;........................................`,

	`                .ckx;'........':c.
             .,:c:c:::oxxocoo::::,',.
            .odc'..:lkkoolllllo;..;d,
            ;c..:o:..;:..',;'.......;.
           ,c..:0Xx::o:.,cllc:,'::,.,c.
           ;c;lkXXXXXXl.;lllll;lXXOo;':c.
         ,dc.oXXXXXXXXl.,lllll;lXXXXx,c0:
         ;Oc.oXXXXXXXXo.':ll:;'oXXXXO;,l'
         'l;;OXXXXXXXXd'.'::'..dXXXXO;,l'
         'l;:0XXXXXXXX0x:...,:o0XXXXk,:x,
         'l;;kXXXXXXKXXXkol;oXXXXXXXO;oNc
        ,c'..ckk0XXXXXXXXXX00XXXXXXX0:;o:.
      .':;..:dd::ooooOXXXXXXXXXXXXXXXo..c;
    .',',:co0XX0kkkxx0XXXXXXXXXXXXXXX0c..;l.
  .:;'..oXXXXXXXXXXXXXXXXXXXXXXXXXXXXXko;';:.
.cdc..:oOXXXXXXXXKXXXXXXXXXXXXXXXXXXXXXXo..oc
:0o...:dxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxo,.:,
cNo........................................;'`,

	`            .cc;.  ...  .;c.
         .,,cc:cc:lxxxl:ccc:;,.
        .lo;...lKKklllookl..cO;
      .cl;.,;'.okl;...'.;,..';:.
     .:o;;dkx,.ll..,cc::,..,'.;:,.
     co..lKKKkokl.':lllo;''ol..;dl.
   .,c;.,xKKKKKKo.':llll;.'oOxo,.cl,.
   cNo..lKKKKKKKo'';llll;;okKKKl..oNc
   cNo..lKKKKKKKko;':c:,'lKKKKKo'.oNc
   cNo..lKKKKKKKKKl.....'dKKKKKxc,l0:
   .c:'.lKKKKKKKKKk;....oKKKKKKo'.oNc
     ,:.,oxOKKKKKKKOxxxxOKKKKKKxc,;ol:.
     ;c..'':oookKKKKKKKKKKKKKKKKKk:.'clc.
   ,dl'.,oxo;'';oxOKKKKKKKKKKKKKKKOxxl::;,,.
  .dOc..lKKKkoooookKKKKKKKKKKKKKKKKKKKxl,;ol.
  cx,';okKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKl..;lc.
  co..:dddddddddddddddddddddddddddddddddl:;''::.
  co..........................................."`,

	`           .ccccccc.
      .,,,;cooolccol;;,,.
     .dOx;..;lllll;..;xOd.
   .cdo,',loOXXXXXkll;';odc.
  ,oo:;c,':oko:cccccc,...ckl.
  ;c.;kXo..::..;c::'.......oc
,dc..oXX0kk0o.':lll;..cxxc.,ld,
kNo.'oXXXXXXo'':lll;..oXXOd;cOd.
KOc;oOXXXXXXo.':lol,..dXXXXl';xc
Ol,:k0XXXXXX0c.,clc'.:0XXXXx,.oc
KOc;dOXXXXXXXl..';'..lXXXXXd..oc
dNo..oXXXXXXXOx:..'lxOXXXXXk,.:; ..
cNo..lXXXXXXXXXOolkXXXXXXXXXkl;..;:.;.
.,;'.,dkkkkk0XXXXXXXXXXXXXXXXXOxxl;,;,;l:.
  ;c.;:''''':doOXXXXXXXXXXXXXXXXXXOdo;';clc.
  ;c.lOdood:'''oXXXXXXXXXXXXXXXXXXXXXk,..;ol.
  ';.:xxxxxocccoxxxxxxxxxxxxxxxxxxxxxxl::'.';;.
  ';........................................;l'`,

	`
        .;:;;,.,;;::,.
     .;':;........'co:.
   .clc;'':cllllc::,.':c.
  .lo;;o:coxdlooollc;',::,,.
.c:'.,cl,.'lc',,;;'......cO;
do;';oxoc::l;;llllc'.';;'.';.
c..ckkkkkkkd,;llllc'.:kkd;.':c.
'.,okkkkkkkkc;llllc,.:kkkdl,cO;
..;xkkkkkkkkc,ccll:,;okkkkk:,cl,
..,dkkkkkkkkc..,;,'ckkkkkkkc;ll.
..'okkkkkkkko,....'okkkkkkkc,:c.
c..ckkkkkkkkkdl;,:okkkkkkkkd,.',';.
d..':lxkkkkkkkkxxkkkkkkkkkkkdoc;,;'..'.,.
o...'';llllldkkkkkkkkkkkkkkkkkkdll;..'cdo.
o..,l;'''''';dkkkkkkkkkkkkkkkkkkkkdlc,..;lc.
o..;lc;;;;;;,,;clllllllllllllllllllllc'..,:c.
o..........................................;'`,

	`
           .,,,,,,,,,.
         .ckKxodooxOOdcc.
      .cclooc'....';;cool.
     .loc;;;;clllllc;;;;;:;,.
   .c:'.,okd;;cdo:::::cl,..oc
  .:o;';okkx;';;,';::;'....,;,.
  co..ckkkkkddk:,cclll;.,c:,:o:.
  co..ckkkkkkkk:,cllll;.:kkd,.':c.
.,:;.,okkkkkkkk:,cclll;.:kkkdl;;o:.
cNo..ckkkkkkkkko,.;llc,.ckkkkkc..oc
,dd;.:kkkkkkkkkx;..;:,.'lkkkkko,.:,
  ;c.ckkkkkkkkkkc.....;ldkkkkkk:.,'
,dc..'okkkkkkkkkxoc;;cxkkkkkkkkc..,;,.
kNo..':lllllldkkkkkkkkkkkkkkkkkdcc,.;l.
KOc,l;''''''';lldkkkkkkkkkkkkkkkkkc..;lc.
xx:':;;;;,.,,...,;;cllllllllllllllc;'.;oo,
cNo.....................................oc`,

	`

                   .ccccccc.
               .ccckNKOOOOkdcc.
            .;;cc:ccccccc:,::::,,.
         .c;:;.,cccllxOOOxlllc,;ol.
        .lkc,coxo:;oOOxooooooo;..:,
      .cdc.,dOOOc..cOd,.',,;'....':c.
      cNx'.lOOOOxlldOl..;lll;.....cO;
     ,do;,:dOOOOOOOOOl'':lll;..:d:.'c,
     co..lOOOOOOOOOOOl'':lll;.'lOd,.cd.
     co.,dOOOOOOOOOOOo,.;llc,.,dOOc..dc
     co..lOOOOOOOOOOOOc.';:,..cOOOl..oc
   .,:;.'::lxOOOOOOOOOo:'...,:oOOOc..dc
   ;Oc..cl'':llxOOOOOOOOdcclxOOOOx,.cd.
  .:;';lxl''''':lldOOOOOOOOOOOOOOc..oc
,dl,.'cooc:::,....,::coooooooooooc'.c:
cNo.................................oc`,

	`


                        .cccccccc.
                  .,,,;;cc:cccccc:;;,.
                .cdxo;..,::cccc::,..;l.
               ,oo:,,:c:cdxxdllll:;,';:,.
             .cl;.,oxxc'.,cc,.',;;'...oNc
             ;Oc..cxxxc'.,c;..;lll;...cO;
           .;;',:ldxxxdoldxc..;lll:'...'c,
           ;c..cxxxxkxxkxxxc'.;lll:'','.cdc.
         .c;.;odxxxxxxxxxxxd;.,cll;.,l:.'dNc
        .:,''ccoxkxxkxxxxxxx:..,:;'.:xc..oNc
      .lc,.'lc':dxxxkxxxxxxxdl,...',lx:..dNc
     .:,',coxoc;;ccccoxxxxxxxxo:::oxxo,.cdc.
  .;':;.'oxxxxxc''''';cccoxxxxxxxxxkxc..oc
,do:'..,:llllll:;;;;;;,..,;:lllllllll;..oc
cNo.....................................oc`,

	`

                              .ccccc.
                         .cc;'coooxkl;.
                     .:c:::c:,;,,,;c;;,.'.
                   .clc,',:,..:xxocc;...c;
                  .c:,';:ox:..:c,,,,,,...cd,
                .c:'.,oxxxxl::l:.;loll;..;ol.
                ;Oc..:xxxxxxxxx:.,llll,....oc
             .,;,',:loxxxxxxxxx:.,llll;.,;.'ld,
            .lo;..:xxxxxxxxxxxx:.'cllc,.:l:'cO;
           .:;...'cxxxxxxxxxxxxol;,::,..cdl;;l'
         .cl;':;'';oxxxxxxxxxxxxx:....,cooc,cO;
     .,,,::;,lxoc:,,:lxxxxxxxxxxxo:,,;lxxl;'oNc
   .cdxo;':lxxxxxxc'';cccccoxxxxxxxxxxxxo,.;lc.
  .loc'.'lxxxxxxxxocc;''''';ccoxxxxxxxxx:..oc
occ'..',:cccccccccccc:;;;;;;;;:ccccccccc,.'c,
Ol;......................................;l'`,

	`
                              ,ddoodd,
                         .cc' ,ooccoo,'cc.
                      .ccldo;....,,...;oxdc.
                   .,,:cc;.''..;lol;;,'..lkl.
                  .dkc';:ccl;..;dl,.''.....oc
                .,lc',cdddddlccld;.,;c::'..,cc:.
                cNo..:ddddddddddd;':clll;,c,';xc
               .lo;,clddddddddddd;':clll;:kc..;'
             .,:;..:ddddddddddddd:';clll;;ll,..
             ;Oc..';:ldddddddddddl,.,c:;';dd;..
           .''',:lc,'cdddddddddddo:,'...'cdd;..
         .cdc';lddd:';lddddddddddddd;.';lddl,..
      .,;::;,cdddddol;;lllllodddddddlcodddd:.'l,
     .dOc..,lddddddddlccc;'';cclddddddddddd;,ll.
   .coc,;::ldddddddddddddl:ccc:ldddddddddlc,ck;
,dl::,..,cccccccccccccccccccccccccccccccc:;':xx,
cNd.........................................;lOc`,
}

var colors = [num_frames]termbox.Attribute{
	// approx colors from original gif
	termbox.Attribute(210), // peach
	termbox.Attribute(222), // orange
	termbox.Attribute(120), // green
	termbox.Attribute(123), // cyan
	termbox.Attribute(111), // blue
	termbox.Attribute(134), // purple
	termbox.Attribute(177), // pink
	termbox.Attribute(207), // fuschia
	termbox.Attribute(206), // magenta
	termbox.Attribute(204), // red
}
