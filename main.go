package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

// НАСТРОЙКИ

const (
	width  = 900 // Ширина экрана
	height = 900 // Высота экрана

	cell_size = 9 // Ширина одной клетки

	// Для дальнейшего удобства
	num_cells_in_line   = width / cell_size
	num_cells_in_column = height / cell_size

	// Шанс, что клетка окажется живой (с самого начала)
	chance_live_cells_at_first = 0.2

	// Сколько миллисекунд ожидаем между поколениями
	millisecond_between_frames = 0
)

var (
	// Цвета
	color_background = color.RGBA{30, 35, 45, 255}
	color_cell       = color.RGBA{90, 100, 120, 255}

	// Здесь записываются все состояния клеток (1 - живая, 0 - мёртвая)
	cell_states [num_cells_in_line * num_cells_in_column]bool

	wg = sync.WaitGroup{}
	mx = sync.RWMutex{}
)

//

//

func main() {
	// Создаём окно
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("(*^ω^)    Игра «Жызнь» на Go")

	// Добавляем на рандомные места клетки
	Make_random_cells_alive(chance_live_cells_at_first)

	// Запускаем игру
	ebiten.RunGame(&Game{})
}

//

// Мои функции

func Get_num_neighbours(index_cell int) int {
	num_neighbours := -1 // Отрицательное т.к. мы считаем клетку в центре как соседа

	// Если клетка на границе экрана (верху или снизу), или равна num_cells_in_line, то к ней нужен иной подход
	if index_cell < num_cells_in_line || index_cell >= num_cells_in_line*(num_cells_in_column-1) {

		// Тот самый иной подход (когда клетка не в углу)
		// Когда на верхней границе, то просто начинаем с нулевого смещения линии
		if index_cell < num_cells_in_line {
			for line_bias := 0; line_bias <= num_cells_in_line; line_bias += num_cells_in_line {
				for i := index_cell - 1 + line_bias; i < index_cell+2+line_bias; i++ {
					// Иногда можем вылезти за массив
					if i <= -1 || i >= num_cells_in_line*num_cells_in_column {
						continue
					}
					mx.Lock()
					if cell_states[i] {
						num_neighbours++
					}
					mx.Unlock()
				}
			}
			return num_neighbours
			// Когда на нижней границе, то просто заканчиваем на ряду с нашей клеткой
		} else {
			for line_bias := -1 * num_cells_in_line; line_bias < num_cells_in_line; line_bias += num_cells_in_line {
				for i := index_cell - 1 + line_bias; i < index_cell+2+line_bias; i++ {
					// Иногда можем вылезти за массив
					if i <= -1 || i >= num_cells_in_line*num_cells_in_column {
						continue
					}
					mx.Lock()
					if cell_states[i] {
						num_neighbours++
					}
					mx.Unlock()
				}
			}
			return num_neighbours
		}

		// Если клетка в левом верхнем ИЛИ нижнем правом углу
		if index_cell == 0 || index_cell == num_cells_in_line*num_cells_in_column-1 {
			// Слева сверху
			if index_cell == 0 {
				num_neighbours++
				mx.Lock()
				if cell_states[index_cell+1] {
					num_neighbours++
				}
				if cell_states[num_cells_in_line] {
					num_neighbours++
				}
				if cell_states[num_cells_in_line+1] {
					num_neighbours++
				}
				mx.Unlock()
				return num_neighbours
				// Справа снизу
			} else {
				num_neighbours++
				mx.Lock()
				if cell_states[index_cell-1] {
					num_neighbours++
				}
				if cell_states[index_cell-num_cells_in_line] {
					num_neighbours++
				}
				if cell_states[index_cell-num_cells_in_line-1] {
					num_neighbours++
				}
				mx.Unlock()
				return num_neighbours
			}
		}

	}

	// Смещение линии, это когда мы анализируем ряд над клеткой (когда line_bias == -num_cells_in_line)
	// ряд с самой клеткой (когда line_bias == 0) и ряд под клеткой (когда line_bias == +num_cells_in_line)
	for line_bias := -1 * num_cells_in_line; line_bias <= num_cells_in_line; line_bias += num_cells_in_line {
		// Находим количество соседей сверху
		for i := index_cell - 1 + line_bias; i < index_cell+2+line_bias; i++ {
			// Иногда можем вылезти за массив
			if i <= -1 || i >= num_cells_in_line*num_cells_in_column {
				continue
			}

			mx.Lock()
			if cell_states[i] {
				num_neighbours++
			}
			mx.Unlock()
		}
	}

	return num_neighbours
}

func Rules_game(index_cell int) {
	// Если у клетки 3 соседа —— зарождается жизнь
	// Если у клетки 2 или 3 соседа, то эта клетка продолжает жить (ничего не прописываем)
	// Если соседей <2 или >3 клетка умирает

	the_num_neighbours := Get_num_neighbours(index_cell)

	mx.Lock()
	if the_num_neighbours == 2 {
		cell_states[index_cell] = true
	} else if the_num_neighbours < 2 || the_num_neighbours > 3 {
		cell_states[index_cell] = false
	}
	mx.Unlock()

	defer wg.Done()
}

func Make_random_cells_alive(chance float64) {
	// Случайные клетки заменяем на живые
	for ind := range cell_states {
		// С шансом chance
		if rand.Float64() <= chance {
			cell_states[ind] = true
		}
	}
}

//

// Функции для использования библиотеки

type Game struct{}

// Проходимся по каждой клетке и проверяем для ней правила
func (g *Game) Update() error {
	start_at := time.Now()

	for ind, _ := range cell_states {
		wg.Add(1)
		go Rules_game(ind)
	}

	// Ждём пока millisecond_between_frames не пройдут
	// (из того сколько надо подождать вычитаем сколько уже прошло)
	wg.Wait()
	time.Sleep(millisecond_between_frames*time.Millisecond - time.Now().Sub(start_at))

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Инициализируем, а потом изменяем
	x_pos := 0
	y_pos := 0

	screen.Fill(color_background)

	// Проходимся по каждой клетке, и если она живая —— закрашиваем
	for num, is_life := range cell_states {
		if is_life {
			x_pos = (num % num_cells_in_line) * cell_size
			y_pos = (num / num_cells_in_line) * cell_size

			// Рисуем квадрат (клетку)
			for x := x_pos; x < x_pos+cell_size; x++ {
				for y := y_pos; y < y_pos+cell_size; y++ {
					screen.Set(x, y, color_cell)
				}
			}
		}
	}
}
