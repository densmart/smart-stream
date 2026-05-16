package com.smartstream.tvclient.ui.details

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.ImageView
import android.widget.TextView
import android.widget.Toast
import androidx.fragment.app.Fragment
import com.bumptech.glide.Glide
import com.smartstream.tvclient.R
import com.smartstream.tvclient.ui.player.PlayerActivity
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Simplified media details fragment
 * Shows: poster, title, duration, and PLAY button
 */
class MediaDetailsFragmentSimple : Fragment() {

    private lateinit var mediaPoster: ImageView
    private lateinit var mediaTitle: TextView
    private lateinit var mediaDuration: TextView
    private lateinit var btnPlay: Button
    private lateinit var btnContinue: Button

    private var mediaId: String? = null
    private var mediaName: String? = null
    private var posterUrl: String? = null
    private var duration: Int = 0
    private var savedPosition: Long = 0L

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Get extras from activity
        activity?.intent?.let { intent ->
            mediaId = intent.getStringExtra(MediaDetailsActivity.EXTRA_MEDIA_ID)
            mediaName = intent.getStringExtra(MediaDetailsActivity.EXTRA_MEDIA_NAME)
            val poster = intent.getStringExtra(MediaDetailsActivity.EXTRA_MEDIA_POSTER)
            duration = intent.getIntExtra(MediaDetailsActivity.EXTRA_MEDIA_DURATION, 0)

            // Build poster URL (same logic as Media.getPosterUrl())
            posterUrl = poster?.let {
                val baseUrl = SharedPrefsManager.getBaseUrl()
                when {
                    it.startsWith("http") -> it
                    it.startsWith("/static/posters/") -> "$baseUrl${it.removePrefix("/")}"
                    else -> "${baseUrl}static/posters/$it"
                }
            }
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_media_details_simple, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        initViews(view)
        displayMediaInfo()

        // Load saved position and update CONTINUE button
        savedPosition = SharedPrefsManager.getMediaPosition(mediaId ?: "")
        updateContinueButton()

        setupPlayButton()
    }

    override fun onStart() {
        super.onStart()
        // Refresh saved position when screen becomes visible
        savedPosition = SharedPrefsManager.getMediaPosition(mediaId ?: "")
        updateContinueButton()
        Log.d("MediaDetails", "onStart: Updated savedPosition = $savedPosition ms")
    }

    override fun onResume() {
        super.onResume()
        // Refresh saved position when returning from player
        savedPosition = SharedPrefsManager.getMediaPosition(mediaId ?: "")
        updateContinueButton()
        Log.d("MediaDetails", "onResume: Updated savedPosition = $savedPosition ms")
    }

    private fun initViews(view: View) {
        mediaPoster = view.findViewById(R.id.media_poster)
        mediaTitle = view.findViewById(R.id.media_title)
        mediaDuration = view.findViewById(R.id.media_duration)
        btnPlay = view.findViewById(R.id.btn_play)
        btnContinue = view.findViewById(R.id.btn_continue)
    }

    private fun displayMediaInfo() {
        // Set title
        mediaTitle.text = mediaName ?: "Unknown Media"

        // Set duration
        if (duration > 0) {
            val durationStr = formatDuration(duration)
            mediaDuration.text = "Duration: $durationStr"
        } else {
            mediaDuration.visibility = View.GONE
        }

        // Load poster
        Log.d("MediaDetails", "posterUrl: $posterUrl")
        if (posterUrl != null) {
            Log.d("MediaDetails", "Loading poster from URL: $posterUrl")
            Glide.with(this)
                .load(posterUrl)
                .error(R.drawable.ic_media_placeholder)
                .placeholder(R.drawable.ic_media_placeholder)
                .into(mediaPoster)
        } else {
            Log.d("MediaDetails", "No poster URL, using placeholder")
            mediaPoster.setImageResource(R.drawable.ic_media_placeholder)
        }
    }

    private fun setupPlayButton() {
        // PLAY button - start from beginning
        btnPlay.setOnClickListener {
            playMedia(startPosition = 0L)
        }

        // CONTINUE button - start from saved position
        btnContinue.setOnClickListener {
            playMedia(startPosition = savedPosition)
        }

        // Request focus for TV navigation
        btnPlay.requestFocus()
    }

    private fun playMedia(startPosition: Long = 0L) {
        mediaId?.let { id ->
            val intent = Intent(requireContext(), PlayerActivity::class.java).apply {
                putExtra(PlayerActivity.EXTRA_MEDIA_ID, id)
                putExtra(PlayerActivity.EXTRA_MEDIA_NAME, mediaName)
                putExtra(PlayerActivity.EXTRA_START_POSITION, startPosition)
            }

            try {
                startActivity(intent)
            } catch (e: Exception) {
                Toast.makeText(
                    requireContext(),
                    "Failed to start player: ${e.message}",
                    Toast.LENGTH_LONG
                ).show()
            }
        } ?: run {
            Toast.makeText(
                requireContext(),
                "Error: Media ID is missing",
                Toast.LENGTH_LONG
            ).show()
        }
    }

    private fun updateContinueButton() {
        if (savedPosition > 0) {
            // Show CONTINUE button with time
            btnContinue.visibility = View.VISIBLE
            val timeStr = formatTimePosition(savedPosition)
            btnContinue.text = "CONTINUE ($timeStr)"
            Log.d("MediaDetails", "CONTINUE button shown with time: $timeStr")
        } else {
            // Hide CONTINUE button
            btnContinue.visibility = View.GONE
            Log.d("MediaDetails", "CONTINUE button hidden (no saved position)")
        }
    }

    private fun formatTimePosition(positionMs: Long): String {
        val totalSeconds = (positionMs / 1000).toInt()
        val hours = totalSeconds / 3600
        val minutes = (totalSeconds % 3600) / 60
        val seconds = totalSeconds % 60

        return String.format("%02d:%02d:%02d", hours, minutes, seconds)
    }

    private fun formatDuration(seconds: Int): String {
        val hours = seconds / 3600
        val minutes = (seconds % 3600) / 60
        val secs = seconds % 60

        return if (hours > 0) {
            String.format("%02d:%02d:%02d", hours, minutes, secs)
        } else {
            String.format("%02d:%02d", minutes, secs)
        }
    }
}
