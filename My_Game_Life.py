#!/usr/local/bin/python
# -*- coding: utf-8 -*-
try:
    import pygame
    from random import randint
    from time import time
    ########################################################## Кастомизируемые

    window_width = 900

    cell_size = 60
    grid_size = 1

    lim_speed = 0.0  # Ограничение скорости на 1 цикл

    percent_lining_cells = 20

    living_cell_color = (45, 45, 65)
    background_color = (200, 200, 200)
    grid_color = (155, 155, 155)

    ########################################################## Не Кастомизируемые

    cells_list = []

    timer = time()

    ##########################################################


    def window():
        """Создаём окно"""
        global wind
        wind = pygame.display.set_mode((window_width, window_width))
        pygame.display.set_caption("(*^ω^)    Игра «Жызнь»")



    def draw():
        """Просто рисуем все клетки и сетку"""
        wind.fill(background_color)

        # Клетки
        for i in cells_list:
            if i[2]:
                pygame.draw.rect(wind, living_cell_color, (i[0], i[1], cell_size, cell_size))

        # Сетка
        for i in range(window_width // cell_size):
            for j in range(window_width // cell_size):
                pygame.draw.rect(wind, grid_color, (i *cell_size, j *cell_size, cell_size +1, cell_size +1), grid_size)

        pygame.display.update()



    def add_cells():
        """Добавляем клетки"""
        global cells_list

        cells_list = []

        for i in range(window_width // cell_size):
            for j in range(window_width // cell_size):
                cells_list.append( [i *cell_size, j *cell_size, True if randint(0,100) <= percent_lining_cells else False] )



    def quit():
        """Узнаём когда закрывать окно"""
        for event in pygame.event.get():
            if event.type == pygame.KEYDOWN:
                # Нажали escape - выходим
                if event.key == pygame.K_ESCAPE:
                    pygame.quit()
                # Нажали R - рестартим
                elif event.key == ord('r'):
                    start_game()
            # Нажали крестик - выходим
            elif event.type == pygame.QUIT:
                pygame.quit()



    def Timer():
        """Отсчитываем определённое время"""
        global timer
        while time() < timer + lim_speed:
            pass
        timer = time()



    def Algorithm_Life():
        """Воплощаем Алгоритм игры Жизнь"""
        # Считаем количество живых соседей
        for i in cells_list:
            amount_neighbors = 0

            # Ищем соседние клетки
            for j in cells_list:
                if j[0] + cell_size == i[0]:
                    if j[1] - cell_size == i[1]:
                        if j[2]:
                            amount_neighbors += 1
                    elif j[1] == i[1]:
                        if j[2]:
                            amount_neighbors += 1
                    elif j[1] + cell_size == i[1]:
                        if j[2]:
                            amount_neighbors += 1

                elif j[0] == i[0]:
                    if j[1] - cell_size == i[1]:
                        if j[2]:
                            amount_neighbors += 1
                    elif j[0] == i[0] and j[1] + cell_size == i[1]:
                        if j[2]:
                            amount_neighbors += 1

                elif j[0] - cell_size == i[0]:
                    if j[1] - cell_size == i[1]:
                        if j[2]:
                            amount_neighbors += 1
                    elif j[1] == i[1]:
                        if j[2]:
                            amount_neighbors += 1
                    elif j[1] + cell_size == i[1]:
                        if j[2]:
                            amount_neighbors += 1


            # Если у клетки есть 3 соседа - воскрешаем
            if amount_neighbors == 3:
                i[2] = True
            # Если у клетки <2 соседей - убиваем
            elif amount_neighbors < 2:
                i[2] = False
            # Если у клетки >3 соседей - убиваем
            elif amount_neighbors > 3:
                i[2] = False


    ##########################################################


    def start_game():
        add_cells()
        while 1:
            draw()
            quit()
            Algorithm_Life()
            Timer()


    window()
    start_game()

except:
    pass