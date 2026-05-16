package com.smartstream.tvclient.ui.main

import android.view.ViewGroup
import androidx.core.content.ContextCompat
import androidx.leanback.widget.ImageCardView
import androidx.leanback.widget.Presenter
import com.bumptech.glide.Glide
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.data.model.Playlist
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Presenter for displaying media and playlist cards in BrowseFragment
 * Uses Leanback's ImageCardView
 */
class CardPresenter : Presenter() {

    companion object {
        private const val CARD_WIDTH = 200
        private const val CARD_HEIGHT = 280
    }

    override fun onCreateViewHolder(parent: ViewGroup): ViewHolder {
        val cardView = ImageCardView(parent.context).apply {
            isFocusable = true
            isFocusableInTouchMode = true
            setBackgroundColor(
                ContextCompat.getColor(context, R.color.tv_card_background)
            )
        }
        return ViewHolder(cardView)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, item: Any) {
        val cardView = viewHolder.view as ImageCardView

        when (item) {
            is Media -> bindMedia(cardView, item)
            is Playlist -> bindPlaylist(cardView, item)
            else -> {
                cardView.setTitleText("Unknown")
                cardView.setContentText("")
            }
        }

        cardView.setMainImageDimensions(CARD_WIDTH, CARD_HEIGHT)
    }

    private fun bindMedia(cardView: ImageCardView, media: Media) {
        cardView.setTitleText(media.name)
        cardView.setContentText(buildMediaContentText(media))

        // Load poster image
        val posterUrl = media.getPosterUrl(SharedPrefsManager.getBaseUrl())
        if (posterUrl != null) {
            Glide.with(cardView.context)
                .load(posterUrl)
                .centerCrop()
                .error(R.drawable.ic_media_placeholder)
                .into(cardView.mainImageView)
        } else {
            cardView.setMainImage(
                ContextCompat.getDrawable(
                    cardView.context,
                    R.drawable.ic_media_placeholder
                )
            )
        }
    }

    private fun bindPlaylist(cardView: ImageCardView, playlist: Playlist) {
        cardView.setTitleText(playlist.name)
        cardView.setContentText(playlist.getTypeDisplayName())

        // Load poster image
        val posterUrl = playlist.getPosterUrl(SharedPrefsManager.getBaseUrl())
        if (posterUrl != null) {
            Glide.with(cardView.context)
                .load(posterUrl)
                .centerCrop()
                .error(R.drawable.ic_playlist_placeholder)
                .into(cardView.mainImageView)
        } else {
            cardView.setMainImage(
                ContextCompat.getDrawable(
                    cardView.context,
                    R.drawable.ic_playlist_placeholder
                )
            )
        }
    }

    private fun buildMediaContentText(media: Media): String {
        val parts = mutableListOf<String>()

        media.getFormattedDuration()?.let { parts.add(it) }
        media.getFormattedSize()?.let { parts.add(it) }

        return parts.joinToString(" • ")
    }

    override fun onUnbindViewHolder(viewHolder: ViewHolder) {
        val cardView = viewHolder.view as ImageCardView
        // Clear image to free memory
        cardView.setBadgeImage(null)
        cardView.setMainImage(null)
    }
}
