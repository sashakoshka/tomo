package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 480, 360))
	window.SetTitle("example label")
	window.Adopt(elements.NewLabelWrapped(text))
	window.OnClose(tomo.Stop)
	window.Show()
}

const text = "Resize the window to see the text wrap:\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Fermentum et sollicitudin ac orci phasellus egestas tellus rutrum. Aliquam vestibulum morbi blandit cursus risus at ultrices mi. Gravida dictum fusce ut placerat. Cursus metus aliquam eleifend mi in nulla posuere. Sit amet nulla facilisi morbi tempus iaculis urna id. Amet volutpat consequat mauris nunc congue nisi vitae. Varius duis at consectetur lorem donec massa sapien faucibus et. Vitae elementum curabitur vitae nunc sed velit dignissim. In hac habitasse platea dictumst quisque sagittis purus. Enim nulla aliquet porttitor lacus luctus accumsan tortor. Lectus magna fringilla urna porttitor rhoncus dolor purus non.\n\nNon pulvinar neque laoreet suspendisse. Viverra adipiscing at in tellus integer. Vulputate dignissim suspendisse in est ante. Purus in mollis nunc sed id semper. In est ante in nibh mauris cursus. Risus pretium quam vulputate dignissim suspendisse in est. Blandit aliquam etiam erat velit scelerisque in dictum. Lectus quam id leo in. Odio tempor orci dapibus ultrices in iaculis. Pharetra sit amet aliquam id. Elit ut aliquam purus sit. Egestas dui id ornare arcu odio ut sem nulla pharetra. Massa tempor nec feugiat nisl pretium fusce id. Dui accumsan sit amet nulla facilisi morbi. A lacus vestibulum sed arcu non odio euismod. Nam libero justo laoreet sit amet cursus. Mattis rhoncus urna neque viverra justo nec. Mauris augue neque gravida in fermentum et sollicitudin ac. Vulputate mi sit amet mauris. Ut sem nulla pharetra diam sit amet."
