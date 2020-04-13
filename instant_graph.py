#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri Apr 10 13:31:34 2020

@author: aurelien
"""

import pandas as pd
import matplotlib.pyplot as plt
from matplotlib import animation
import numpy as np
import sys
import os

def main(args):
    
    # Filename
    dirname = args[0]
    
    dataframe = pd.read_csv(os.path.join(dirname, "positions.csv"), delimiter=";")
    nb_times = len(dataframe['time'].unique())
    times = np.array(dataframe['time'].unique()).reshape(nb_times,)
    nb_agents = len(dataframe['id'].unique())
    
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
        
    fig = plt.figure()
    fig.patch.set_facecolor((1, 1, 1))
    ax = plt.axes(xlim=(0, times[-1]), ylim=(0, nb_agents))
    plt.tick_params(axis = "x", which = "both", bottom = False, top = False)
    plt.xticks([], " ")
    plt.title("Change over time")
    plt.ylabel("Number")
    plt.xlabel("Time")
    
    ax.fill_between(times, cum3, color='black')
    ax.fill_between(times, cum2, color='green')
    ax.fill_between(times, cum1,color=(0.8,0.8,0.8))
    ax.fill_between(times, sick, color='red')
    
    plt.show()

if __name__ == "__main__":
    main(sys.argv[1:])
