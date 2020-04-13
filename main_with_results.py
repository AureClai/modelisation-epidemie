#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Wed Apr  8 17:32:00 2020

@author: aurelien
"""

import pandas as pd
import matplotlib.pyplot as plt
from matplotlib import animation
import numpy as np
import json

save_anim = False

dataframe = pd.read_csv('simulation_positions.csv', delimiter=";")
times = dataframe['time'].unique()
filename = "rendered.mp4"

def dot(vect1, vect2):
    return vect1[0]*vect2[0] + vect1[1]*vect2[1]

def norm(vect):
    return np.sqrt(dot(vect, vect))

def get_wall_angle(wall):
    if wall['end']['x']-wall['start']['x'] != 0:
        angle = np.arctan(np.abs(wall['end']['y']-wall['start']['y'])/np.abs(wall['end']['x']-wall['start']['x']))
        if wall['end']['x']-wall['start']['x'] <0 and wall['end']['y']-wall['start']['y']>0:
            # second quarter
            angle += np.pi/2
        elif wall['end']['x']-wall['start']['x'] <0 and wall['end']['y']-wall['start']['y']<=0:
            # third quarter
            angle += np.pi
        elif wall['end']['x']-wall['start']['x'] >0 and wall['end']['y']-wall['start']['y']<0:
            #fourth quarter
            angle += 3*np.pi/2
            
    else:
        if wall['end']['y']-wall['start']['y'] > 0:
            angle = np.pi/2
        else:
            angle = -np.pi/2
    return angle*180/np.pi

with open('settings.json') as settingsFile:
    settings = json.load(settingsFile)
    
    walls = settings['walls']
    window_size_x = settings['window_size_x']
    window_size_y = settings['window_size_y']
    radius = settings['agents_radius']

colors = {
    0 : 'grey',
    1 : 'red',
    2 : 'green',
    3 : 'black'
    }



subDf = dataframe[dataframe.time == times[0]]

fig = plt.figure()
ax = plt.axes(xlim=(0, window_size_x), ylim=(0, window_size_y))
ax.set_aspect('equal')
for wall in walls:
    rect1_patch = plt.Rectangle((wall['start']['x'], wall['start']['y']), width=norm([wall['end']['x']-wall['start']['x'], wall['end']['y']-wall['start']['y']]), height=wall['radius'], angle=get_wall_angle(wall), color='black')
    rect2_patch = plt.Rectangle((wall['end']['x'], wall['end']['y']), width=norm([wall['end']['x']-wall['start']['x'], wall['end']['y']-wall['start']['y']]), height=wall['radius'], angle=180+get_wall_angle(wall), color='black')
    circ1_patch = plt.Circle((wall['start']['x'], wall['start']['y']), radius=wall['radius'], fc='black')
    circ2_patch = plt.Circle((wall['end']['x'], wall['end']['y']), radius=wall['radius'], fc='black')
    ax.add_patch(rect1_patch)
    ax.add_patch(rect2_patch)
    ax.add_patch(circ1_patch)
    ax.add_patch(circ2_patch)
    

patches = {}
for index, row in subDf.iterrows():
    patches[row['id']] = plt.Circle((row['x'], row['y']),radius=radius, fc=colors[row['state']])
    ax.add_patch(patches[row['id']])
    
    
def animate(i, times, dataframe, patches):
    print(times[i])
    subDf = dataframe[dataframe.time == times[i]]
    for index, row in subDf.iterrows():
        patches[row['id']].center = (row['x'], row['y'])
        patches[row['id']].set_color(colors[row['state']])
    

anim = animation.FuncAnimation(fig, animate, 
                               frames=len(times), 
                               interval=1000*(times[1]- times[0]),
                               fargs=(times, dataframe, patches))

if save_anim:
    mywriter = animation.FFMpegWriter(fps=int(1/(times[1]-times[0])))
    anim.save(filename,writer=mywriter)
    print("Animation saved to " + filename)