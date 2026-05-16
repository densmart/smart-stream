package com.smartstream.tvclient

import android.app.Application

class SmartStreamApp : Application() {

    override fun onCreate() {
        super.onCreate()
        instance = this
    }

    companion object {
        lateinit var instance: SmartStreamApp
            private set
    }
}
