/**
 * 
 */

joint.shapes.basic.Generic.define('crl.Element',
    {
        attrs: {
            rect: {
                width: 300
            },
            '.bounding-rect': {
                stroke: "black",
                magnet: true,
                'stroke-width': 2,
                fill: "#ffffff",
                height: 40,
                width: 200,
                transform: "translate(0,0)"
            },
            '.image': {
                'type':'image',
                'ref-x': 1.0,
                'ref-y': 1.0,
                ref: ".bounding-rect",
                width: 16,
                height: 16,
                'xlink:href': ""
            },
            '.abstractions-text': {
                ref: ".bounding-rect",
                'ref-y': 1.0,
                'ref-x': 1.0 + 18,
                'text-anchor': "left",
                'y-alignment': "top",
                'font-weight': "normal",
                'font-style': "italic",
                fill: "black",
                'font-size': 10,
                'font-family': "Go,  Helvetica, Ariel, sans-serif",
                'text':"defaultAbstractionsText"
            },
            '.labelText': {
                ref: ".bounding-rect",
                'ref-y': 4.0 + 16,
                'ref-x': 3.0,
                'text-anchor': "left",
                'y-alignment': "top",
                'font-weight': "bold",
                fill: "black",
                'font-size': 12,
                'font-family': "Go,  Helvetica, Ariel, sans-serif",
                'text':""
            }
        },
        crlJointID: "",
        name:  "labelDefault", 
        abstractions: "defaultAbstractions",
        icon:"",
        displayLabelYOffset:0.0
    },
    {
        markup: "<g class=\"rotatable\">" +
        "<g class=\"scalable\">" +
        "<rect class=\"bounding-rect\"/>" +
        "</g>" +
        "<image class=\"image\"/>" +
        "<text class=\"abstractions-text\"/>" + 
        "<text class=\"labelText\" />" +
        "</g>",

        initialize: function () {
            this.on('change:name change:abstractions change:icon change:displayLabelYOffset', function () {
               this.updateRectangles();
               this.trigger('crl-update');
            }, this);
            this.on("change:position", crlOnChangePosition);
            this.updateRectangles();
            joint.shapes.basic.Generic.prototype.initialize.apply(this, arguments);
        },

        updateRectangles: function () {
            var attrs = this.get('attrs');
            var boundingRectAttr = attrs['.bounding-rect'];
            attrs['.labelText'].text = this.get("name");
            attrs['.image']['xlink:href'] = this.get("icon");
            attrs['.abstractions-text'].text = this.get("abstractions");
            attrs['.labelText']['ref-y'] = this.get("displayLabelYOffset");
            this.resize(boundingRectAttr.width , boundingRectAttr.height);
        }
    });

    joint.shapes.crl.ElementView = joint.dia.ElementView.extend({}, {
        initialize: function() {
            joint.dia.ElementView.prototype.initialize.apply(this, arguments);
            this.listenTo(this.model, 'crl-update', function() {
                this.update();
                this.resize();
            });
        }
    });

    joint.dia.Link.define('crl.OwnerPointer', {
        attrs: {
            line: {
                connection: true,
                stroke: '#333333',
                strokeWidth: 2,
                strokeLinejoin: 'round',
                targetMarker: {
                    "type": "path",
                    "d": "M 10 -5 0 0 10 5 20 0 z"
                }
            },
            wrapper: {
                connection: true,
                strokeWidth: 10,
                strokeLinejoin: 'round'
            }
        },
        crlJointID: ""
    }, {
        markup: [{
            tagName: 'path',
            selector: 'wrapper',
            attributes: {
                'fill': 'none',
                'cursor': 'pointer',
                'stroke': 'transparent'
            }
        }, {
            tagName: 'path',
            selector: 'line',
            attributes: {
                'fill': 'none',
                'pointer-events': 'none'
            }
        }]
    });
    
    joint.dia.Link.define('crl.ElementPointer', {
        attrs: {
            line: {
                connection: true,
                stroke: '#333333',
                strokeWidth: 2,
                strokeLinejoin: 'round',
                targetMarker: {
                    "type": "path",
                    "d": "M 10 -5 0 0 10 5 z"
                }
            },
            wrapper: {
                connection: true,
                strokeWidth: 10,
                strokeLinejoin: 'round'
            }
        },
        crlJointID: ""
    }, {
        markup: [{
            tagName: 'path',
            selector: 'wrapper',
            attributes: {
                'fill': 'none',
                'cursor': 'pointer',
                'stroke': 'transparent'
            }
        }, {
            tagName: 'path',
            selector: 'line',
            attributes: {
                'fill': 'none',
                'pointer-events': 'none'
            }
        }]
    });
    
    joint.dia.Link.define('crl.AbstractPointer', {
        attrs: {
            line: {
                connection: true,
                stroke: '#333333',
                strokeWidth: 2,
                strokeLinejoin: 'round',
                sourceMarker: {
                    "type": "path",
                    "fill": "white",
                    "d": "M 0 -8 15 0 0 8 z"
                }
            },
            wrapper: {
                connection: true,
                strokeWidth: 10,
                strokeLinejoin: 'round'
            }
        },
        crlJointID: ""
    }, {
        markup: [{
            tagName: 'path',
            selector: 'wrapper',
            attributes: {
                'fill': 'none',
                'cursor': 'pointer',
                'stroke': 'transparent'
            }
        }, {
            tagName: 'path',
            selector: 'line',
            attributes: {
                'fill': 'none',
                'pointer-events': 'none'
            }
        }]
    });
    
    joint.dia.Link.define('crl.RefinedPointer', {
        attrs: {
            line: {
                connection: true,
                stroke: '#333333',
                strokeWidth: 2,
                strokeLinejoin: 'round',
                sourceMarker: {
                    "type": "path",
                    "fill": "white",
                    "d": "M 15 -8 0 0 15 8 z"
                },
                targetMarker: {
                    "type":"path",
                    "d": "M 10 -5 0 0 10 5 z"
                }
            },
            wrapper: {
                connection: true,
                strokeWidth: 10,
                strokeLinejoin: 'round'
            }
        },
        crlJointID: ""
    }, {
        markup: [{
            tagName: 'path',
            selector: 'wrapper',
            attributes: {
                'fill': 'none',
                'cursor': 'pointer',
                'stroke': 'transparent'
            }
        }, {
            tagName: 'path',
            selector: 'line',
            attributes: {
                'fill': 'none',
                'pointer-events': 'none'
            }
        }]
    });
    



    joint.dia.Link.define('crl.ReferenceLink', {
        attrs: {
            line: {
                connection: true,
                stroke: '#333333',
                strokeWidth: 2,
                strokeLinejoin: 'round',
                targetMarker: {
                    'type': 'path',
                    'd': 'M 10 -5 0 0 10 5 z'
                },
                sourceMarker: {
                    "type": "path",
                    "d": "M 10 -5 0 0 10 5 20 0 z"
                }
            },
            wrapper: {
                connection: true,
                strokeWidth: 10,
                strokeLinejoin: 'round'
            }
        }
    }, {
        markup: [{
            tagName: 'path',
            selector: 'wrapper',
            attributes: {
                'fill': 'none',
                'cursor': 'pointer',
                'stroke': 'transparent'
            }
        }, {
            tagName: 'path',
            selector: 'line',
            attributes: {
                'fill': 'none',
                'pointer-events': 'none'
            }
        }]
    });
    
    joint.dia.Link.define('crl.RefinementLink', {
        attrs: {
            line: {
                connection: true,
                stroke: 'black',
                strokeWidth: 2,
                strokeLinejoin: 'round',
                targetMarker: {
                    'type': 'path',
                    "fill":"white",
                    'd': 'M 15 -7 0 0 15 7 z'
                }
            },
            wrapper: {
                connection: true,
                strokeWidth: 10,
                strokeLinejoin: 'round'
            }
        }
    }, {
        markup: [{
            tagName: 'path',
            selector: 'wrapper',
            attributes: {
                'fill': 'none',
                'cursor': 'pointer',
                'stroke': 'transparent'
            }
        }, {
            tagName: 'path',
            selector: 'line',
            attributes: {
                'fill': 'none',
                'pointer-events': 'none'
            }
        }]
    });
