#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri Apr 10 14:02:14 2020

@author: aurelien
"""

import pandas as pd
import matplotlib.pyplot as plt
from matplotlib import animation
import matplotlib as mpl
import numpy as np
import json
import sys
import os

mpl.rcParams["savefig.facecolor"] = ((1, 1, 1))
mpl.rcParams.update({'font.size': 12})



# Functions
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
    
    # Script 
def main(args):
    
     # Filename
    dirname = args[0]
    
    dataframe = pd.read_csv(os.path.join(dirname, "positions.csv"), delimiter=";")
    nb_times = len(dataframe['time'].unique())
    times = dataframe['time'].unique()
    nb_agents = len(dataframe['id'].unique())
    
    with open(os.path.join(dirname,'settings.json')) as settingsFile:
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
    
    window_size_x = 30
    window_size_y = 30
    radius = 0.2
    
    healthy = np.zeros((nb_times,))
    dead = np.zeros((nb_times,))
    sick = np.zeros((nb_times,))
    recovered = np.zeros((nb_times,))
    # healthy + sick
    cum1 = np.zeros((nb_times,))
    # healthy + sick + recovered
    cum2 = np.zeros((nb_times,))
    # all
    cum3 = np.zeros((nb_times,))
    
    for i in range(len(times)):
        subDf = dataframe[dataframe.time == times[i]]
        healthy[i] = len(subDf[subDf['state']==0])
        sick[i] = len(subDf[subDf['state']==1])
        recovered[i] = len(subDf[subDf['state']==2])
        dead[i] = len(subDf[subDf['state']==3])
        cum1[i] = healthy[i] + sick[i]
        cum2[i] = healthy[i] + sick[i] + recovered[i]
        cum3[i] = healthy[i] + sick[i] + recovered[i] + dead[i]
    
    subDf = dataframe[dataframe.time == times[0]]
    
    f, (a0, a1) = plt.subplots(2, 1, gridspec_kw={'height_ratios': [1, 3]})
    f.set_size_inches([6.4, 9.36])
    a0.plot([0],[0])
    a0.set_xlim(0, times[-1])
    a0.set_ylim(0, nb_agents)
    a0.tick_params(axis = "x", which = "both", bottom = False, top = False)
    a0.set_xticklabels([], fontdict=None, minor=False)
    a0.set_yticklabels([], fontdict=None, minor=False)
    a0.set_xlabel('Temps')
    a0.set_ylabel('Nombre')
    a1.set_xlim(0, window_size_x)
    a1.set_ylim(0, window_size_y)
    a1.set_aspect('equal')
    a1.tick_params(axis = "x", which = "both", bottom = False, top = False)
    a1.tick_params(axis = "y", which = "both", left = False, right = False)
    a1.set_xticklabels([], fontdict=None, minor=False)
    a1.set_yticklabels([], fontdict=None, minor=False)
    
    patches = {}
    for index, row in subDf.iterrows():
        patches[row['id']] = plt.Circle((row['x'], row['y']),radius=radius, fc=colors[row['state']])
        a1.add_patch(patches[row['id']])
        
    for wall in walls:
        rect1_patch = plt.Rectangle((wall['start']['x'], wall['start']['y']), width=norm([wall['end']['x']-wall['start']['x'], wall['end']['y']-wall['start']['y']]), height=wall['radius'], angle=get_wall_angle(wall), color='black')
        rect2_patch = plt.Rectangle((wall['end']['x'], wall['end']['y']), width=norm([wall['end']['x']-wall['start']['x'], wall['end']['y']-wall['start']['y']]), height=wall['radius'], angle=180+get_wall_angle(wall), color='black')
        circ1_patch = plt.Circle((wall['start']['x'], wall['start']['y']), radius=wall['radius'], fc='black')
        circ2_patch = plt.Circle((wall['end']['x'], wall['end']['y']), radius=wall['radius'], fc='black')
        a1.add_patch(rect1_patch)
        a1.add_patch(rect2_patch)
        a1.add_patch(circ1_patch)
        a1.add_patch(circ2_patch)
        
    def animate(i, a0, times, sick, cum1, cum2, cum3, dataframe, patches ):
        print(times[i])
        # first graph
        a0.collections.clear()
        a0.fill_between(times[:i], cum3[:i], color='black')
        a0.fill_between(times[:i], cum2[:i], color='green')
        a0.fill_between(times[:i], cum1[:i],color=(0.8,0.8,0.8))
        a0.fill_between(times[:i], sick[:i], color='red')
        # second graph
        subDf = dataframe[dataframe.time == times[i]]
        for index, row in subDf.iterrows():
            patches[row['id']].center = (row['x'], row['y'])
            patches[row['id']].set_color(colors[row['state']])
        
    anim = animation.FuncAnimation(f, animate, 
                                   frames=int(len(times)*1), 
                                   interval=1000*(times[1]- times[0]),
                                   fargs=(a0, times, sick, cum1, cum2, cum3, dataframe, patches))
        
    mywriter = animation.FFMpegWriter(fps=60)
    anim.save(os.path.join(dirname , 'video.mp4'),writer=mywriter, dpi=400)
    
if __name__ == "__main__":
    main(sys.argv[1:])
