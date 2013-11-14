import QtQuick 2.0

Item {
    id: block

    property int type: 0

    Image {
        id: img

        anchors.fill: parent
        source: {
            if (type == 0)
                return "redStone.png"
            else if (type == 1)
                return "blueStone.png"
            else
                return "greenStone.png"
        }
    }
}
