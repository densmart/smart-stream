package com.smartstream.tvclient.ui.details

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Adapter for vertical list of media items
 */
class MediaListAdapter(
    private val playlistId: String?,
    private val onMediaClick: (Media) -> Unit
) : RecyclerView.Adapter<MediaListAdapter.MediaViewHolder>() {

    private val mediaList = mutableListOf<Media>()

    fun submitList(list: List<Media>) {
        mediaList.clear()
        // Sort by order field
        mediaList.addAll(list.sortedBy { it.order })
        notifyDataSetChanged()
    }

    fun refreshList() {
        notifyDataSetChanged()
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): MediaViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_media_list_row, parent, false)
        return MediaViewHolder(view, playlistId, onMediaClick)
    }

    override fun onBindViewHolder(holder: MediaViewHolder, position: Int) {
        holder.bind(mediaList[position])
    }

    override fun getItemCount(): Int = mediaList.size

    class MediaViewHolder(
        itemView: View,
        private val playlistId: String?,
        private val onMediaClick: (Media) -> Unit
    ) : RecyclerView.ViewHolder(itemView) {

        private val poster: ImageView = itemView.findViewById(R.id.media_poster)
        private val name: TextView = itemView.findViewById(R.id.media_name)
        private val duration: TextView = itemView.findViewById(R.id.media_duration)
        private val continueWatching: TextView = itemView.findViewById(R.id.media_continue)

        fun bind(media: Media) {
            name.text = media.name
            duration.text = media.getFormattedDuration() ?: "N/A"

            // Check if this media is the current episode in playlist
            if (playlistId != null) {
                val currentMediaId = SharedPrefsManager.getPlaylistCurrentMedia(playlistId)
                val savedPosition = SharedPrefsManager.getPlaylistPosition(playlistId)

                if (currentMediaId == media.id && savedPosition > 1000) {
                    // This is the current episode with saved position
                    val timeText = formatTime(savedPosition)
                    continueWatching.text = itemView.context.getString(R.string.continue_watching, timeText)
                    continueWatching.visibility = View.VISIBLE
                } else {
                    continueWatching.visibility = View.GONE
                }
            } else {
                // For single media (not in playlist), check individual media position
                val savedPosition = SharedPrefsManager.getMediaPosition(media.id)
                if (savedPosition > 1000) {
                    val timeText = formatTime(savedPosition)
                    continueWatching.text = itemView.context.getString(R.string.continue_watching, timeText)
                    continueWatching.visibility = View.VISIBLE
                } else {
                    continueWatching.visibility = View.GONE
                }
            }

            // Load poster
            val posterUrl = media.getPosterUrl(SharedPrefsManager.getBaseUrl())
            if (posterUrl != null) {
                Glide.with(itemView.context)
                    .load(posterUrl)
                    .centerCrop()
                    .error(R.drawable.ic_media_placeholder)
                    .into(poster)
            } else {
                poster.setImageResource(R.drawable.ic_media_placeholder)
            }

            // Click on the whole row
            itemView.setOnClickListener {
                onMediaClick(media)
            }
        }

        private fun formatTime(milliseconds: Long): String {
            val totalSeconds = milliseconds / 1000
            val hours = totalSeconds / 3600
            val minutes = (totalSeconds % 3600) / 60
            val seconds = totalSeconds % 60

            return if (hours > 0) {
                String.format("%02d:%02d:%02d", hours, minutes, seconds)
            } else {
                String.format("%02d:%02d", minutes, seconds)
            }
        }
    }
}
