package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/sebastianappelberg/disk/pkg/clean"
	"github.com/sebastianappelberg/disk/pkg/storage"
	"github.com/spf13/cobra"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

const minTableWidth = 67

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Delete  key.Binding
	Exclude key.Binding
	Exit    key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Delete, k.Exclude, k.Exit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Delete}}
}

type model struct {
	table          table.Model
	help           help.Model
	keyMap         KeyMap
	windowWidth    int
	dialogWidth    int
	dialogHeight   int
	tableWidth     int
	totalReclaimed int64
	total          int64
	dir            string
	cleanableFiles []clean.CleanableFile
	inProgressWg   *sync.WaitGroup
}

func (m model) Init() tea.Cmd {
	//return m.spinner.Tick
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		height := msg.Height - 5
		m.windowWidth = msg.Width
		m.table.SetHeight(height)
		m.dialogWidth = msg.Width - m.tableWidth - 14
		m.dialogHeight = height + 3
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "up", "k", "down", "j":
			// TODO: Use to make rows more compact.
			m.dir = filepath.Dir(m.cleanableFiles[m.table.Cursor()].Path)
		case "q", "ctrl+c":
			return m, tea.Quit
		case "e", "enter":
			cursor := m.table.Cursor()
			if len(m.table.Rows()) > 0 && cursor < len(m.table.Rows()) {
				file := m.cleanableFiles[cursor]
				m.total -= file.Size
				m.asyncAction(func() {
					file.Exclude()
				})
				slices.Delete(m.cleanableFiles, cursor, cursor+1)
				m.table.SetRows(slices.Delete(m.table.Rows(), cursor, cursor+1))
			}
		case "w", "backspace":
			cursor := m.table.Cursor()
			if len(m.table.Rows()) > 0 && cursor < len(m.table.Rows()) {
				file := m.cleanableFiles[cursor]
				m.total += file.Size
				m.asyncAction(func() {
					err := file.Remove()
					if err != nil {
						log.Printf("error putting %q in the trash: %v", file.Path, err)
					}
				})
				slices.Delete(m.cleanableFiles, cursor, cursor+1)
				m.table.SetRows(slices.Delete(m.table.Rows(), cursor, cursor+1))
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) asyncAction(action func()) {
	m.inProgressWg.Add(1)
	go func() {
		defer m.inProgressWg.Done()
		action()
	}()
}

func (m model) View() string {
	if m.windowWidth <= 185 {
		// If the window is too narrow then skip rendering the help dialog.
		return lipgloss.JoinVertical(lipgloss.Top,
			m.tableView(),
			m.bottomView(),
		)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Top,
			m.tableView(),
			m.bottomView(),
		),
		m.dialogView(),
	)
}

func (m model) tableView() string {
	return baseStyle.Render(m.table.View())
}

func (m model) bottomView() string {
	return baseStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Right,
			" "+m.help.ShortHelpView(m.keyMap.ShortHelp()),
			m.reclaimedSpaceView(),
		))
}

func (m model) dialogView() string {
	in := `# disk clean

This is a list of pesky space hoggers that **disk clean** _thinks_ can be removed safely.
Don't worry if you accidentally delete something, deleting from this list means moving it to the recycling bin.
Speaking of which, to complete the cleaning you'll have to empty your bin.

Examples of files and folders it will suggest:
- Clutter in the form caches, dependency folders, build folders, etc. above a given size and age.
- Steam games you haven't played in a while.
- Movies and TV shows that are easy to get a hold of even if you delete them.

If you exclude a file it will be excluded for all future runs of the **disk clean** command.
To reset your excluded files, delete the "$HOME/.disk/user_config_cache" file. 
`
	r, err := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.dialogWidth),
	)
	if err != nil {
		log.Fatal(err)
	}
	out, err := r.Render(in)
	if err != nil {
		log.Fatal(err)
	}
	return baseStyle.
		Height(m.dialogHeight).Render(out)
}

var rightStyle = lipgloss.NewStyle().Align(lipgloss.Right).PaddingRight(1)

func (m model) reclaimedSpaceView() string {
	return rightStyle.
		Width(m.tableWidth - 67).
		Render(fmt.Sprintf("Reclaimed space: %s/%s", storage.FormatSize(m.totalReclaimed), storage.FormatSize(m.total)))
}

func NewCmdClean() *cobra.Command {
	var minSize int
	var minAge int
	var maxPlaytime int

	var cmd = &cobra.Command{
		Use:   "clean <path>",
		Short: "Show a recommendation of folders and files to delete to save space.",
		Long: `Shows a TUI with a table of files and folders that have been marked as candidates for deletion.

Examples of files and folders it will suggest:
- Clutter in the form caches, dependency folders, build folders, etc. above a given size and age.
- Steam games you haven't played in a while.
- Movies and TV shows that are easy to get a hold of even if you delete them. 
`,
		Run: func(cmd *cobra.Command, args []string) {
			// To be nice on the user's CPU this command will only use 1/2 of the available CPUs.
			runtime.GOMAXPROCS(int(math.Ceil(float64(runtime.NumCPU() / 2))))

			root := args[0]

			var rows []table.Row
			longestPath := 0
			total := int64(0)

			cleanableFiles := clean.Clean(clean.Args{
				Root:        root,
				MinAge:      minAge,
				MinSize:     minSize,
				MaxPlaytime: maxPlaytime,
			})

			for _, file := range cleanableFiles {
				path := strings.TrimPrefix(file.Path, root)
				if len(path) > longestPath {
					longestPath = len(path)
				}
				rows = append(rows, table.Row{
					path,
					storage.FormatSize(file.Size),
					file.ModTime.Format(time.DateTime),
				})
				total += file.Size
			}

			sizeColWidth := 8
			lastUsedColWidth := 20
			if longestPath < minTableWidth {
				longestPath = minTableWidth
			}
			columns := []table.Column{
				{Title: "Folder", Width: longestPath},
				{Title: "Size", Width: sizeColWidth},
				{Title: "Last Used", Width: lastUsedColWidth},
			}

			keyMap := KeyMap{
				Up: key.NewBinding(
					key.WithKeys("up", "k"),
					key.WithHelp("↑/k", "up"),
				),
				Down: key.NewBinding(
					key.WithKeys("down", "j"),
					key.WithHelp("↓/j", "down"),
				),
				Delete: key.NewBinding(
					key.WithKeys("w", "backspace"),
					key.WithHelp("w/backspace", "delete"),
				),
				Exclude: key.NewBinding(
					key.WithKeys("e", "enter"),
					key.WithHelp("e/enter", "exclude"),
				),
				Exit: key.NewBinding(
					key.WithKeys("q", "ctrl+c"),
					key.WithHelp("q/ctrl+c", "quit"),
				),
			}

			t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithFocused(true))
			s := table.DefaultStyles()
			s.Header = s.Header.
				BorderForeground(lipgloss.Color("240")).
				Bold(true)
			s.Selected = s.Selected.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)
			t.SetStyles(s)

			tableWidth := longestPath + sizeColWidth + lastUsedColWidth

			m := model{
				table:          t,
				help:           help.New(),
				keyMap:         keyMap,
				tableWidth:     tableWidth,
				total:          total,
				cleanableFiles: cleanableFiles,
				inProgressWg:   &sync.WaitGroup{},
			}

			_, err := tea.NewProgram(m, tea.WithMouseCellMotion(), tea.WithAltScreen()).Run()
			if err != nil {
				log.Fatal(err)
			}

			// TODO: Run spinner.
			m.inProgressWg.Wait()
			// TODO: Print deletion report. x files deleted, x GB reclaimed.
		},
	}

	cmd.Flags().IntVarP(&minSize, "min-size", "s", 50, "Minimum size of files to include in analysis results specified in megabytes.")
	cmd.Flags().IntVarP(&minAge, "min-age", "a", 90, "Minimum age of files to include in analysis results specified in days.")
	cmd.Flags().IntVarP(&maxPlaytime, "max-playtime", "p", 20, "Maximum playtime of games to include in analysis results specified in hours.")

	return cmd
}
