package com.smartstream.tvclient.ui.player

import android.os.Bundle
import android.util.Log
import android.widget.Toast
import androidx.fragment.app.FragmentActivity
import com.google.android.exoplayer2.ExoPlayer
import com.google.android.exoplayer2.MediaItem
import com.google.android.exoplayer2.PlaybackException
import com.google.android.exoplayer2.Player
import com.google.android.exoplayer2.source.DefaultMediaSourceFactory
import com.google.android.exoplayer2.ui.PlayerView
import com.google.android.exoplayer2.upstream.DataSource
import com.google.android.exoplayer2.upstream.DefaultDataSource
import com.google.android.exoplayer2.upstream.DefaultHttpDataSource
import com.smartstream.tvclient.R
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Video player activity for TV.
 * Uses ExoPlayer to stream media from server.
 */
class PlayerActivity : FragmentActivity() {

    private var player: ExoPlayer? = null
    private lateinit var playerView: PlayerView

    private var mediaId: String? = null
    private var mediaName: String? = null
    private var startPosition: Long = 0L

    companion object {
        private const val TAG = "PlayerActivity"
        const val EXTRA_MEDIA_ID = "media_id"
        const val EXTRA_MEDIA_NAME = "media_name"
        const val EXTRA_START_POSITION = "start_position"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        Log.d(TAG, "onCreate: =============== PlayerActivity started ===============")
        super.onCreate(savedInstanceState)

        try {
            setContentView(R.layout.activity_player)
            Log.d(TAG, "onCreate: Layout inflated successfully")

            // Get media info from intent
            mediaId = intent.getStringExtra(EXTRA_MEDIA_ID)
            mediaName = intent.getStringExtra(EXTRA_MEDIA_NAME)
            startPosition = intent.getLongExtra(EXTRA_START_POSITION, 0L)

            Log.d(TAG, "onCreate: mediaId=$mediaId, mediaName=$mediaName, startPosition=$startPosition")

            if (mediaId == null) {
                Log.e(TAG, "onCreate: mediaId is null, finishing activity")
                Toast.makeText(this, "Error: Media ID is missing", Toast.LENGTH_LONG).show()
                finish()
                return
            }

            // Check if token exists
            val token = SharedPrefsManager.getToken()
            if (token == null) {
                Log.e(TAG, "onCreate: JWT token is null!")
                Toast.makeText(this, "Error: Authentication required", Toast.LENGTH_LONG).show()
                finish()
                return
            }
            Log.d(TAG, "onCreate: JWT token exists (length=${token.length})")

            // Check server configuration
            val baseUrl = SharedPrefsManager.getBaseUrl()
            Log.d(TAG, "onCreate: Base URL = $baseUrl")

            playerView = findViewById(R.id.player_view)
            Log.d(TAG, "onCreate: PlayerView found")

            initializePlayer()
        } catch (e: Exception) {
            Log.e(TAG, "onCreate: Exception occurred!", e)
            Toast.makeText(this, "Error: ${e.message}", Toast.LENGTH_LONG).show()
            finish()
        }
    }

    private fun initializePlayer() {
        Log.d(TAG, "initializePlayer: Starting player initialization")

        try {
            // Create DataSource.Factory with JWT authentication
            val dataSourceFactory = createDataSourceFactory()
            Log.d(TAG, "initializePlayer: DataSourceFactory created")

            // Build streaming URL
            val streamUrl = buildStreamUrl(mediaId!!)
            Log.d(TAG, "initializePlayer: Stream URL = $streamUrl")

            // Show URL to user for debugging
            Toast.makeText(this, "Loading: $streamUrl", Toast.LENGTH_SHORT).show()

            // Create player
            Log.d(TAG, "initializePlayer: Creating ExoPlayer instance")
            player = ExoPlayer.Builder(this)
                .setMediaSourceFactory(DefaultMediaSourceFactory(dataSourceFactory))
                .build()
                .also { exoPlayer ->
                    Log.d(TAG, "initializePlayer: ExoPlayer instance created")

                    playerView.player = exoPlayer
                    Log.d(TAG, "initializePlayer: PlayerView attached to ExoPlayer")

                    val mediaItem = MediaItem.fromUri(streamUrl)
                    Log.d(TAG, "initializePlayer: MediaItem created from URI")

                    exoPlayer.setMediaItem(mediaItem)
                    Log.d(TAG, "initializePlayer: MediaItem set on player")

                    exoPlayer.prepare()
                    Log.d(TAG, "initializePlayer: Player prepare() called")

                    // Seek to start position if provided
                    if (startPosition > 0) {
                        exoPlayer.seekTo(startPosition)
                        Log.d(TAG, "initializePlayer: Seeked to position $startPosition ms")
                    }

                    exoPlayer.playWhenReady = true
                    Log.d(TAG, "initializePlayer: playWhenReady set to true")

                    // Add player state listener
                    exoPlayer.addListener(object : Player.Listener {
                        override fun onPlaybackStateChanged(playbackState: Int) {
                            val stateString = when (playbackState) {
                                Player.STATE_IDLE -> "IDLE"
                                Player.STATE_BUFFERING -> "BUFFERING"
                                Player.STATE_READY -> "READY"
                                Player.STATE_ENDED -> "ENDED"
                                else -> "UNKNOWN"
                            }
                            Log.d(TAG, "onPlaybackStateChanged: $stateString")

                            if (playbackState == Player.STATE_READY) {
                                Toast.makeText(this@PlayerActivity, "Playback started", Toast.LENGTH_SHORT).show()
                            } else if (playbackState == Player.STATE_ENDED) {
                                // Video finished, clear saved position and close player
                                mediaId?.let { SharedPrefsManager.clearMediaPosition(it) }
                                Log.d(TAG, "onPlaybackStateChanged: Video ended, cleared saved position")

                                // Close player after a short delay
                                playerView.postDelayed({
                                    finish()
                                }, 1000)
                            }
                        }

                        override fun onIsPlayingChanged(isPlaying: Boolean) {
                            Log.d(TAG, "onIsPlayingChanged: isPlaying=$isPlaying")
                        }

                        override fun onPlayerError(error: PlaybackException) {
                            Log.e(TAG, "onPlayerError: ${error.message}", error)
                            Log.e(TAG, "onPlayerError: errorCode=${error.errorCode}")
                            Log.e(TAG, "onPlayerError: cause=${error.cause?.message}")
                            handlePlayerError(error)
                        }
                    })

                    Log.d(TAG, "initializePlayer: Player listener added")
                }

            Log.d(TAG, "initializePlayer: Player initialization completed successfully")
        } catch (e: Exception) {
            Log.e(TAG, "initializePlayer: Exception during initialization!", e)
            Toast.makeText(this, "Player initialization error: ${e.message}", Toast.LENGTH_LONG).show()
            finish()
        }
    }

    private fun createDataSourceFactory(): DataSource.Factory {
        Log.d(TAG, "createDataSourceFactory: Creating data source factory")

        val jwtToken = SharedPrefsManager.getToken()
        Log.d(TAG, "createDataSourceFactory: JWT token present = ${jwtToken != null}")

        if (jwtToken != null) {
            Log.d(TAG, "createDataSourceFactory: Token length = ${jwtToken.length}")
            Log.d(TAG, "createDataSourceFactory: Token preview = ${jwtToken.take(20)}...")
        }

        // Create HTTP data source with JWT token
        val httpDataSourceFactory = DefaultHttpDataSource.Factory()
            .setUserAgent("SmartStreamTVClient/1.0")
            .setConnectTimeoutMs(60000) // 60 seconds
            .setReadTimeoutMs(60000) // 60 seconds
            .setAllowCrossProtocolRedirects(true)

        Log.d(TAG, "createDataSourceFactory: HTTP data source factory configured")
        Log.d(TAG, "createDataSourceFactory: Connect timeout = 60s, Read timeout = 60s")

        // Add Authorization header if token exists
        if (jwtToken != null) {
            httpDataSourceFactory.setDefaultRequestProperties(
                mapOf("Authorization" to "Bearer $jwtToken")
            )
            Log.d(TAG, "createDataSourceFactory: Authorization header added to requests")
        } else {
            Log.w(TAG, "createDataSourceFactory: WARNING - No JWT token found!")
        }

        val dataSourceFactory = DefaultDataSource.Factory(this, httpDataSourceFactory)
        Log.d(TAG, "createDataSourceFactory: Data source factory created successfully")

        return dataSourceFactory
    }

    private fun buildStreamUrl(mediaId: String): String {
        val baseUrl = SharedPrefsManager.getBaseUrl()
        return "${baseUrl}media/$mediaId/stream/"
    }

    private fun handlePlayerError(error: PlaybackException) {
        Log.e(TAG, "handlePlayerError: Processing error...")
        Log.e(TAG, "handlePlayerError: Error type = ${error.javaClass.simpleName}")

        val errorMessage = when (error.errorCode) {
            PlaybackException.ERROR_CODE_IO_NETWORK_CONNECTION_FAILED -> {
                Log.e(TAG, "handlePlayerError: Network connection failed")
                "Network connection failed"
            }
            PlaybackException.ERROR_CODE_IO_NETWORK_CONNECTION_TIMEOUT -> {
                Log.e(TAG, "handlePlayerError: Network timeout")
                "Network timeout"
            }
            PlaybackException.ERROR_CODE_IO_UNSPECIFIED -> {
                Log.e(TAG, "handlePlayerError: IO error")
                "IO error: ${error.cause?.message ?: "Unknown"}"
            }
            PlaybackException.ERROR_CODE_IO_BAD_HTTP_STATUS -> {
                Log.e(TAG, "handlePlayerError: Bad HTTP status")
                "Bad HTTP status: ${error.cause?.message ?: "Unknown"}"
            }
            PlaybackException.ERROR_CODE_BEHIND_LIVE_WINDOW -> {
                Log.e(TAG, "handlePlayerError: Behind live window")
                "Behind live window"
            }
            PlaybackException.ERROR_CODE_TIMEOUT -> {
                Log.e(TAG, "handlePlayerError: Timeout")
                "Playback timeout"
            }
            else -> {
                Log.e(TAG, "handlePlayerError: Other error (code=${error.errorCode})")
                "Playback error: ${error.message ?: "Unknown error"}"
            }
        }

        Log.e(TAG, "handlePlayerError: Showing error to user: $errorMessage")
        Toast.makeText(this, errorMessage, Toast.LENGTH_LONG).show()

        // Don't finish immediately, give user time to see the error
        playerView.postDelayed({
            finish()
        }, 3000)
    }

    override fun onPause() {
        super.onPause()
        Log.d(TAG, "onPause: Pausing playback and saving position")
        player?.playWhenReady = false

        // Save current position when pausing (so it's available when returning to details screen)
        player?.let { exoPlayer ->
            val currentPosition = exoPlayer.currentPosition
            val duration = exoPlayer.duration

            Log.d(TAG, "onPause: currentPosition=$currentPosition, duration=$duration")

            // Only save if position is valid and not at the very end (within 1 second)
            if (currentPosition > 1000 && (duration <= 0 || currentPosition < duration - 1000)) {
                mediaId?.let { id ->
                    SharedPrefsManager.saveMediaPosition(id, currentPosition)
                    Log.d(TAG, "onPause: Saved position $currentPosition ms for media $id")
                }
            } else {
                Log.d(TAG, "onPause: Position NOT saved (currentPosition=$currentPosition, duration=$duration)")
            }
        }
    }

    override fun onResume() {
        super.onResume()
        Log.d(TAG, "onResume: Resuming playback")
        player?.playWhenReady = true
    }

    override fun onDestroy() {
        super.onDestroy()
        Log.d(TAG, "onDestroy: Activity is being destroyed")
        releasePlayer()
    }

    private fun releasePlayer() {
        Log.d(TAG, "releasePlayer: Releasing player resources")
        player?.release()
        player = null
    }
}
