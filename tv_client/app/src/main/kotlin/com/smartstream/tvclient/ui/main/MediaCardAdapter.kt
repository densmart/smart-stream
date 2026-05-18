package com.smartstream.tvclient.ui.main

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.data.model.Playlist
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Adapter for displaying media and playlist cards in RecyclerView
 */
class MediaCardAdapter(
    private val onItemClick: (Any) -> Unit,
    private val onItemFocus: (Any) -> Unit,
    private val onFocusLost: () -> Unit = {}
) : RecyclerView.Adapter<MediaCardAdapter.CardViewHolder>() {

    private val items = mutableListOf<Any>()

    fun submitList(newItems: List<Any>) {
        items.clear()
        items.addAll(newItems)
        notifyDataSetChanged()
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): CardViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_media_card, parent, false)
        return CardViewHolder(view, onItemClick, onItemFocus, onFocusLost)
    }

    override fun onBindViewHolder(holder: CardViewHolder, position: Int) {
        holder.bind(items[position])
    }

    override fun getItemCount(): Int = items.size

    class CardViewHolder(
        itemView: View,
        private val onItemClick: (Any) -> Unit,
        private val onItemFocus: (Any) -> Unit,
        private val onFocusLost: () -> Unit
    ) : RecyclerView.ViewHolder(itemView) {

        private val poster: ImageView = itemView.findViewById(R.id.card_poster)
        private val title: TextView = itemView.findViewById(R.id.card_title)
        private val info: TextView = itemView.findViewById(R.id.card_info)

        fun bind(item: Any) {
            when (item) {
                is Media -> bindMedia(item)
                is Playlist -> bindPlaylist(item)
            }

            // Click listener
            itemView.setOnClickListener {
                onItemClick(item)
            }

            // Focus listener
            itemView.onFocusChangeListener = View.OnFocusChangeListener { _, hasFocus ->
                if (hasFocus) {
                    onItemFocus(item)
                    // Enable marquee scrolling for title
                    title.isSelected = true
                } else {
                    // Disable marquee scrolling when focus is lost
                    title.isSelected = false
                    // Notify that focus was lost
                    onFocusLost()
                }
            }
        }

        private fun bindMedia(media: Media) {
            title.text = media.name
            info.text = buildMediaInfo(media)

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
        }

        private fun bindPlaylist(playlist: Playlist) {
            title.text = playlist.name
            info.text = playlist.getTypeDisplayName(itemView.context)

            // Load poster
            val posterUrl = playlist.getPosterUrl(SharedPrefsManager.getBaseUrl())
            if (posterUrl != null) {
                Glide.with(itemView.context)
                    .load(posterUrl)
                    .centerCrop()
                    .error(R.drawable.ic_playlist_placeholder)
                    .into(poster)
            } else {
                poster.setImageResource(R.drawable.ic_playlist_placeholder)
            }
        }

        private fun buildMediaInfo(media: Media): String {
            val parts = mutableListOf<String>()
            media.getFormattedDuration()?.let { parts.add(it) }
            media.getFormattedSize()?.let { parts.add(it) }
            return parts.joinToString(" • ")
        }
    }
}
