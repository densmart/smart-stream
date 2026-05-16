package com.smartstream.tvclient.ui.details

import androidx.leanback.widget.AbstractDetailsDescriptionPresenter

/**
 * Presenter for media details description.
 * Displays title and description in details overview.
 */
class MediaDetailsDescriptionPresenter : AbstractDetailsDescriptionPresenter() {

    override fun onBindDescription(vh: ViewHolder, item: Any) {
        if (item is String) {
            vh.body.text = item
        }
    }
}
