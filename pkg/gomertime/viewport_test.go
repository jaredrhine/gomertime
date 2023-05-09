package gomertime

import "testing"

func TestTextViewportCalc(t *testing.T) {
	cases := []struct {
		description     string
		label           string
		worldX          int
		worldY          int
		viewportX       int
		viewportY       int
		width           int
		height          int
		expectedShow    bool
		expectedScreenX int
		expectedScreenY int
		expectedIcon    string
	}{
		{
			description:     "position at origin is shown",
			label:           "entity",
			worldX:          0,
			worldY:          0,
			viewportX:       0,
			viewportY:       0,
			width:           50,
			height:          20,
			expectedShow:    true,
			expectedScreenX: 1,
			expectedScreenY: 3,
			expectedIcon:    "E",
		},
		{
			description:     "position to right along the x-axis is shown",
			label:           "entity",
			worldX:          3,
			worldY:          0,
			viewportX:       0,
			viewportY:       0,
			width:           50,
			height:          20,
			expectedShow:    true,
			expectedScreenX: 4,
			expectedScreenY: 3,
			expectedIcon:    "E",
		},
		{
			description:     "position to right along the x-axis too far is not shown",
			label:           "entity",
			worldX:          51,
			worldY:          0,
			viewportX:       0,
			viewportY:       0,
			width:           50,
			height:          20,
			expectedShow:    false,
			expectedScreenX: 0,
			expectedScreenY: 0,
			expectedIcon:    "",
		},
		{
			description:     "position down along the y-axis is shown",
			label:           "entity",
			worldX:          3,
			worldY:          -10,
			viewportX:       0,
			viewportY:       0,
			width:           50,
			height:          20,
			expectedShow:    true,
			expectedScreenX: 4,
			expectedScreenY: 13,
			expectedIcon:    "E",
		},
		{
			description:     "position down along the y-axis too far is not shown",
			label:           "entity",
			worldX:          3,
			worldY:          -21,
			viewportX:       0,
			viewportY:       0,
			width:           50,
			height:          20,
			expectedShow:    false,
			expectedScreenX: 0,
			expectedScreenY: 0,
			expectedIcon:    "",
		},
		{
			description:     "position up the y-axis is not shown",
			label:           "entity",
			worldX:          0,
			worldY:          3,
			viewportX:       0,
			viewportY:       0,
			width:           50,
			height:          20,
			expectedShow:    false,
			expectedScreenX: 0,
			expectedScreenY: 0,
			expectedIcon:    "",
		},
		{
			description:     "position when viewport moved is shown",
			label:           "entity",
			worldX:          13,
			worldY:          -3,
			viewportX:       5,
			viewportY:       3,
			width:           50,
			height:          20,
			expectedShow:    true,
			expectedScreenX: 9,
			expectedScreenY: 9,
			expectedIcon:    "E",
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			show, x, y, icon := TextViewportCalc(tt.label, tt.worldX, tt.worldY, tt.viewportX, tt.viewportY, tt.width, tt.height, 2, 3)
			if show != tt.expectedShow || x != tt.expectedScreenX || y != tt.expectedScreenY || icon != tt.expectedIcon {
				t.Errorf("show got %t wanted %t, worldX got %d wanted %d, worldY got %d wanted %d", show, tt.expectedShow, x, tt.expectedScreenX, y, tt.expectedScreenY)
			}
		})
	}
}
